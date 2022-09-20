package unittest

import (
	"context"
	"errors"
	"github.com/arvians-id/apriori/http/request"
	"github.com/arvians-id/apriori/http/response"
	"github.com/arvians-id/apriori/model"
	repository "github.com/arvians-id/apriori/repository/mock"
	"github.com/arvians-id/apriori/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

var userRepository = &repository.UserRepositoryMock{
	Mock: mock.Mock{},
}
var userService = service.UserServiceImpl{
	UserRepository: userRepository,
}

var userRequest = model.User{
	IdUser:    1,
	Role:      1,
	Name:      "Widdy",
	Email:     "widdy@gmail.com",
	Address:   "Jl Bhayangkara",
	Phone:     "082299921720",
	Password:  "Rahasia123.",
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
}

func TestFindByAll(t *testing.T) {
	var userRequests []*model.User
	userRequests = append(userRequests, &userRequest, &model.User{
		IdUser:    2,
		Role:      1,
		Name:      "Arfiansyah",
		Email:     "arfiansyah@gmail.com",
		Address:   "Jl Bhayangkara",
		Phone:     "082299921720",
		Password:  "Rahasia123.",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	t.Run("when the list user is not found, then return error", func(t *testing.T) {
		test := userRepository.Mock.On("FindAll", mock.Anything).Return(nil, nil)

		ctx := context.Background()
		users, err := userService.FindAll(ctx)
		assert.NotNil(t, err)
		assert.Nil(t, users)
		assert.Equal(t, err, errors.New(response.ErrorNotFound))
		test.Unset()
	})

	t.Run("when the list user is found, then return list users", func(t *testing.T) {
		test := userRepository.Mock.On("FindAll", mock.Anything).Return(userRequests, nil)

		ctx := context.Background()
		users, err := userService.FindAll(ctx)
		assert.NotNil(t, users)
		assert.Nil(t, err)
		assert.Equal(t, userRequests[0].Name, users[0].Name)
		assert.Equal(t, 2, len(users))
		test.Unset()
	})
}

func TestFindById(t *testing.T) {
	t.Run("when the user is not found, then return error", func(t *testing.T) {
		test := userRepository.Mock.On("FindById", mock.Anything, 1).Return(nil, nil)

		ctx := context.Background()
		user, err := userService.FindById(ctx, 1)
		assert.NotNil(t, err)
		assert.Nil(t, user)
		assert.Equal(t, err, errors.New(response.ErrorNotFound))
		test.Unset()
	})

	t.Run("when the user is found, then return user details", func(t *testing.T) {
		test := userRepository.Mock.On("FindById", mock.Anything, 1).Return(&userRequest, nil)

		ctx := context.Background()
		user, err := userService.FindById(ctx, 1)
		assert.NotNil(t, user)
		assert.Nil(t, err)
		assert.Equal(t, userRequest.Name, user.Name)
		test.Unset()
	})
}

func TestFindByEmail(t *testing.T) {
	t.Run("when the user is not found, then return error", func(t *testing.T) {
		test := userRepository.Mock.On("FindByEmail", mock.Anything, "widdy@gmail.com").Return(nil, nil)

		ctx := context.Background()
		user, err := userService.FindByEmail(ctx, &request.GetUserCredentialRequest{
			Email:    userRequest.Email,
			Password: "Rahasia123.",
		})
		assert.NotNil(t, err)
		assert.Nil(t, user)
		assert.Equal(t, err, errors.New(response.ErrorNotFound))
		test.Unset()
	})

	t.Run("when the user is found, then return user details", func(t *testing.T) {
		passwordHashed, _ := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
		userRequest.Password = string(passwordHashed)
		test := userRepository.Mock.On("FindByEmail", mock.Anything, "widdy@gmail.com").Return(&userRequest, nil)

		ctx := context.Background()
		user, err := userService.FindByEmail(ctx, &request.GetUserCredentialRequest{
			Email:    userRequest.Email,
			Password: "Rahasia123.",
		})
		assert.NotNil(t, user)
		assert.Nil(t, err)
		assert.Equal(t, userRequest.Name, user.Name)
		test.Unset()
	})
}

func TestCreate(t *testing.T) {
	test := userRepository.Mock.On("Create", mock.Anything).Return(&userRequest, nil)

	ctx := context.Background()
	user, err := userService.Create(ctx, &request.CreateUserRequest{
		Name:     userRequest.Name,
		Email:    userRequest.Email,
		Address:  userRequest.Address,
		Phone:    userRequest.Phone,
		Password: userRequest.Password,
	})

	assert.NotNil(t, user)
	assert.Nil(t, err)
	assert.Equal(t, userRequest.Name, user.Name)
	test.Unset()
}

func TestUpdate(t *testing.T) {
	t.Run("when the user is not found, then return error", func(t *testing.T) {
		test := userRepository.Mock.On("FindById", mock.Anything, 1).Return(nil, nil)
		test2 := userRepository.Mock.On("Update", mock.Anything).Return(nil, nil)

		ctx := context.Background()
		user, err := userService.Update(ctx, &request.UpdateUserRequest{
			IdUser:  1,
			Role:    2,
			Name:    "Arfiansyah",
			Email:   "widdy@ummi.ac.id",
			Address: "Jl Bhayangkara",
			Phone:   "082299921720",
		})
		assert.NotNil(t, err)
		assert.Nil(t, user)
		assert.Equal(t, err, errors.New(response.ErrorNotFound))
		test.Unset()
		test2.Unset()
	})

	t.Run("when the user is found, then update it", func(t *testing.T) {
		userRequestUpdate := model.User{
			IdUser:  1,
			Role:    2,
			Name:    "Arfiansyah",
			Email:   "widdy@ummi.ac.id",
			Address: "Jl Bhayangkara",
			Phone:   "082299921720",
		}
		test := userRepository.Mock.On("FindById", mock.Anything, 1).Return(&userRequest, nil)
		test2 := userRepository.Mock.On("Update", mock.Anything).Return(&userRequestUpdate, nil)

		ctx := context.Background()
		user, err := userService.Update(ctx, &request.UpdateUserRequest{
			IdUser:  userRequestUpdate.IdUser,
			Role:    userRequestUpdate.Role,
			Name:    userRequestUpdate.Name,
			Email:   userRequestUpdate.Email,
			Address: userRequestUpdate.Address,
			Phone:   userRequestUpdate.Phone,
		})
		assert.NotNil(t, user)
		assert.Nil(t, err)
		assert.NotEqual(t, "Widdy", userRequest.Name)
		assert.Equal(t, "widdy@ummi.ac.id", userRequest.Email)
		test.Unset()
		test2.Unset()
	})
}

func TestDelete(t *testing.T) {
	t.Run("when the user is not found, then return error", func(t *testing.T) {
		test := userRepository.Mock.On("FindById", mock.Anything, 1).Return(nil, nil)
		test2 := userRepository.Mock.On("Delete", mock.Anything, 1).Return(nil)

		ctx := context.Background()
		err := userService.Delete(ctx, 1)
		assert.NotNil(t, err)
		assert.Equal(t, err, errors.New(response.ErrorNotFound))
		test.Unset()
		test2.Unset()
	})

	t.Run("when unsuccessful deleting user, then return error", func(t *testing.T) {
		test := userRepository.Mock.On("FindById", mock.Anything, 1).Return(&userRequest, nil)
		test2 := userRepository.Mock.On("Delete", mock.Anything, 1).Return(errors.New("err"))

		ctx := context.Background()
		err := userService.Delete(ctx, 1)
		assert.NotNil(t, err)
		assert.Equal(t, err, errors.New("something went wrong"))
		test.Unset()
		test2.Unset()
	})

	t.Run("when the user is found, then delete it", func(t *testing.T) {
		test := userRepository.Mock.On("FindById", mock.Anything, 1).Return(&userRequest, nil)
		test2 := userRepository.Mock.On("Delete", mock.Anything, 1).Return(nil)

		ctx := context.Background()
		err := userService.Delete(ctx, 1)
		assert.Nil(t, err)
		test.Unset()
		test2.Unset()
	})
}
