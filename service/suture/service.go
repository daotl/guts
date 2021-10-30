package suturesrv

import (
	"context"
	"errors"

	"github.com/thejerf/suture/v4"
	"go.uber.org/zap"

	"github.com/daotl/go-log/v2"
	gsync "github.com/daotl/guts/sync"
)

var closedchan = make(chan struct{})

func init() {
	close(closedchan)
}

var (
	ErrRunFnMustBeSpecified = errors.New("runFn must be specified")
	ErrAlreadyRunning       = errors.New("already running")
	ErrNotRunning           = errors.New("not running")
)

// Status of a Service
type Status uint8

const (
	// Not started or stopped
	StatusStopped Status = iota
	// Running but not yet ready to serve
	StatusRunning
	// Running and ready to serve
	StatusReady
	// Stopping
	StatusStopping
)

/*
Service is a Suture-based service interface.

Service's lifecycle graph:
                                Serve()
                 StatusStopped ─────────► StatusRunning
                 │                                 │ │
runFn() returned │ ┌───────────────────────────────┘ │ ready() called in RunFunc
                 │ ▼  Stop() or context canceled     ▼
                 StatusStopping ◄───────── StatusReady
*/
type Service interface {
	suture.Service

	// Serve is called by a Supervisor to start the service.
	// Returns ErrAlreadyRunning if it's already running.
	// See: https://pkg.go.dev/github.com/thejerf/suture/v4#hdr-Serve_Method
	Serve(ctx context.Context) error

	// Status returns the status of the service.
	Status() Status

	// Ready returns a channel to notify that the service is ready to serve.
	// Note that the service may be stopped and restarted before the returned channel is closed.
	Ready() chan struct{}

	// Stop the service if it's running, wait on the returned channel for stopping to finish.
	// Returns ErrNotRunning if it's not running.
	Stop() (chan struct{}, error)
}

// RunFunc is called when service is started either by a Supervisor or manually by calling Service.Serve.
// The service should execute within the goroutine that this is called in, that is, it should not spawn a
// "worker" goroutine. If this function either returns error or panics, the Supervisor will call it again.
// `ready` should be called when the service is ready to serve.
// Note that a service may stop before it's ready.
type RunFunc = func(ctx context.Context, ready func()) error

// BaseService is the base implementation of a Service which can be embedded in actual services.
type BaseService struct {
	Logger log.StandardLogger

	cancel  context.CancelFunc
	mtx     gsync.Mutex
	status  Status
	ready   chan struct{}
	stopped chan struct{}

	// runFn will run when the service is started.
	runFn RunFunc
}

var _ = Service(&BaseService{})

// NewBaseService creates a new BaseService.
func NewBaseService(runFn RunFunc, logger log.StandardLogger,
) (*BaseService, error) {
	if runFn == nil {
		return nil, ErrRunFnMustBeSpecified
	}
	if logger == nil {
		logger = zap.NewNop().Sugar()
	}

	return &BaseService{
		Logger: logger,
		runFn:  runFn,
		ready:  make(chan struct{}),
	}, nil
}

func (s *BaseService) Status() Status {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return s.status
}

func (s *BaseService) Serve(ctx context.Context) error {
	s.mtx.Lock()
	{
		if s.status != StatusStopped {
			s.mtx.Unlock()
			return ErrAlreadyRunning
		}
		s.status = StatusRunning
		// Don't s.ready here.
		// Reset s.ready when and only when the service become ready, so that the waiting goroutines
		// can continue waiting on the same channel until it's ready even the service restarts.
		//s.ready = make(chan struct{})
		s.stopped = make(chan struct{})
		ctx, s.cancel = context.WithCancel(ctx)
	}
	s.mtx.Unlock()

	go func() {
		<-ctx.Done()
		s.mtx.Lock()
		defer s.mtx.Unlock()
		if s.status == StatusRunning || s.status == StatusReady {
			s.status = StatusStopping
		}
	}()

	s.Logger.Debug("Service started running")
	defer s.Logger.Debug("Service stopped")
	err := s.runFn(ctx, func() {
		s.mtx.Lock()
		defer s.mtx.Unlock()
		if s.status == StatusRunning {
			s.status = StatusReady
			close(s.ready)
			// Reset s.ready when and only when the service become ready, so that the waiting goroutines
			// can continue waiting on the same channel until it's ready even the service restarts.
			s.ready = make(chan struct{})
		}
	})

	s.mtx.Lock()
	{
		s.status = StatusStopped
		close(s.stopped)
	}
	s.mtx.Unlock()

	return err
}

func (s *BaseService) Ready() chan struct{} {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	// s.ready is reset when the service become ready, so we have to return a closed channel here
	if s.status == StatusReady {
		return closedchan
	}
	return s.ready
}

func (s *BaseService) Stop() (chan struct{}, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.status == StatusStopped {
		return nil, ErrNotRunning
	}
	if s.status != StatusStopping {
		s.status = StatusStopping
		s.cancel()
	}

	return s.stopped, nil
}

// BaseSupervisor is the base implementation of a Suture Supervisor which can be embedded in actual supervisors.
type BaseSupervisor struct {
	*suture.Supervisor

	Logger log.StandardLogger

	cancel context.CancelFunc
}

// Serve starts the supervisor. You should call this on the top-level supervisor, but nothing else.
func (s *BaseSupervisor) Serve(ctx context.Context) error {
	s.Logger.Debug("Supervisor started running")
	defer s.Logger.Debug("Supervisor stopped")
	ctx, s.cancel = context.WithCancel(ctx)
	return s.Supervisor.Serve(ctx)
}

// ServeBackground starts running a supervisor in its own goroutine.
// When this method returns, the supervisor is guaranteed to be in a running state.
// The returned one-buffered channel receives the error returned by Serve.
func (s *BaseSupervisor) ServeBackground(ctx context.Context) <-chan error {
	s.Logger.Debug("Supervisor started running in the background")
	ctx, s.cancel = context.WithCancel(ctx)
	return s.Supervisor.ServeBackground(ctx)
}

// NewBaseSupervisor creates a new BaseSupervisor.
func NewBaseSupervisor(logger log.StandardLogger,
) (*BaseSupervisor, error) {
	if logger == nil {
		logger = zap.NewNop().Sugar()
	}

	return &BaseSupervisor{
		Logger: logger,
	}, nil
}
