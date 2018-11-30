// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package follow

import (
	"sync"

	"github.com/dwarvesf/yggdrasil/services/networks/model"
	"github.com/satori/go.uuid"
)

var (
	lockServiceMockFindAll  sync.RWMutex
	lockServiceMockFollow   sync.RWMutex
	lockServiceMockSave     sync.RWMutex
	lockServiceMockUnFollow sync.RWMutex
)

// ServiceMock is a mock implementation of Service.
//
//     func TestSomethingThatUsesService(t *testing.T) {
//
//         // make and configure a mocked Service
//         mockedService := &ServiceMock{
//             FindAllFunc: func(q *Query) ([]model.Follow, error) {
// 	               panic("TODO: mock out the FindAll method")
//             },
//             FollowFunc: func(fromUser uuid.UUID, toUser uuid.UUID) error {
// 	               panic("TODO: mock out the Follow method")
//             },
//             SaveFunc: func(r *model.Follow) error {
// 	               panic("TODO: mock out the Save method")
//             },
//             UnFollowFunc: func(fromUser uuid.UUID, toUser uuid.UUID) error {
// 	               panic("TODO: mock out the UnFollow method")
//             },
//         }
//
//         // TODO: use mockedService in code that requires Service
//         //       and then make assertions.
//
//     }
type ServiceMock struct {
	// FindAllFunc mocks the FindAll method.
	FindAllFunc func(q *Query) ([]model.Follow, error)

	// FollowFunc mocks the Follow method.
	FollowFunc func(fromUser uuid.UUID, toUser uuid.UUID) error

	// SaveFunc mocks the Save method.
	SaveFunc func(r *model.Follow) error

	// UnFollowFunc mocks the UnFollow method.
	UnFollowFunc func(fromUser uuid.UUID, toUser uuid.UUID) error

	// calls tracks calls to the methods.
	calls struct {
		// FindAll holds details about calls to the FindAll method.
		FindAll []struct {
			// Q is the q argument value.
			Q *Query
		}
		// Follow holds details about calls to the Follow method.
		Follow []struct {
			// FromUser is the fromUser argument value.
			FromUser uuid.UUID
			// ToUser is the toUser argument value.
			ToUser uuid.UUID
		}
		// Save holds details about calls to the Save method.
		Save []struct {
			// R is the r argument value.
			R *model.Follow
		}
		// UnFollow holds details about calls to the UnFollow method.
		UnFollow []struct {
			// FromUser is the fromUser argument value.
			FromUser uuid.UUID
			// ToUser is the toUser argument value.
			ToUser uuid.UUID
		}
	}
}

// FindAll calls FindAllFunc.
func (mock *ServiceMock) FindAll(q *Query) ([]model.Follow, error) {
	if mock.FindAllFunc == nil {
		panic("ServiceMock.FindAllFunc: method is nil but Service.FindAll was just called")
	}
	callInfo := struct {
		Q *Query
	}{
		Q: q,
	}
	lockServiceMockFindAll.Lock()
	mock.calls.FindAll = append(mock.calls.FindAll, callInfo)
	lockServiceMockFindAll.Unlock()
	return mock.FindAllFunc(q)
}

// FindAllCalls gets all the calls that were made to FindAll.
// Check the length with:
//     len(mockedService.FindAllCalls())
func (mock *ServiceMock) FindAllCalls() []struct {
	Q *Query
} {
	var calls []struct {
		Q *Query
	}
	lockServiceMockFindAll.RLock()
	calls = mock.calls.FindAll
	lockServiceMockFindAll.RUnlock()
	return calls
}

// Follow calls FollowFunc.
func (mock *ServiceMock) Follow(fromUser uuid.UUID, toUser uuid.UUID) error {
	if mock.FollowFunc == nil {
		panic("ServiceMock.FollowFunc: method is nil but Service.Follow was just called")
	}
	callInfo := struct {
		FromUser uuid.UUID
		ToUser   uuid.UUID
	}{
		FromUser: fromUser,
		ToUser:   toUser,
	}
	lockServiceMockFollow.Lock()
	mock.calls.Follow = append(mock.calls.Follow, callInfo)
	lockServiceMockFollow.Unlock()
	return mock.FollowFunc(fromUser, toUser)
}

// FollowCalls gets all the calls that were made to Follow.
// Check the length with:
//     len(mockedService.FollowCalls())
func (mock *ServiceMock) FollowCalls() []struct {
	FromUser uuid.UUID
	ToUser   uuid.UUID
} {
	var calls []struct {
		FromUser uuid.UUID
		ToUser   uuid.UUID
	}
	lockServiceMockFollow.RLock()
	calls = mock.calls.Follow
	lockServiceMockFollow.RUnlock()
	return calls
}

// Save calls SaveFunc.
func (mock *ServiceMock) Save(r *model.Follow) error {
	if mock.SaveFunc == nil {
		panic("ServiceMock.SaveFunc: method is nil but Service.Save was just called")
	}
	callInfo := struct {
		R *model.Follow
	}{
		R: r,
	}
	lockServiceMockSave.Lock()
	mock.calls.Save = append(mock.calls.Save, callInfo)
	lockServiceMockSave.Unlock()
	return mock.SaveFunc(r)
}

// SaveCalls gets all the calls that were made to Save.
// Check the length with:
//     len(mockedService.SaveCalls())
func (mock *ServiceMock) SaveCalls() []struct {
	R *model.Follow
} {
	var calls []struct {
		R *model.Follow
	}
	lockServiceMockSave.RLock()
	calls = mock.calls.Save
	lockServiceMockSave.RUnlock()
	return calls
}

// UnFollow calls UnFollowFunc.
func (mock *ServiceMock) UnFollow(fromUser uuid.UUID, toUser uuid.UUID) error {
	if mock.UnFollowFunc == nil {
		panic("ServiceMock.UnFollowFunc: method is nil but Service.UnFollow was just called")
	}
	callInfo := struct {
		FromUser uuid.UUID
		ToUser   uuid.UUID
	}{
		FromUser: fromUser,
		ToUser:   toUser,
	}
	lockServiceMockUnFollow.Lock()
	mock.calls.UnFollow = append(mock.calls.UnFollow, callInfo)
	lockServiceMockUnFollow.Unlock()
	return mock.UnFollowFunc(fromUser, toUser)
}

// UnFollowCalls gets all the calls that were made to UnFollow.
// Check the length with:
//     len(mockedService.UnFollowCalls())
func (mock *ServiceMock) UnFollowCalls() []struct {
	FromUser uuid.UUID
	ToUser   uuid.UUID
} {
	var calls []struct {
		FromUser uuid.UUID
		ToUser   uuid.UUID
	}
	lockServiceMockUnFollow.RLock()
	calls = mock.calls.UnFollow
	lockServiceMockUnFollow.RUnlock()
	return calls
}