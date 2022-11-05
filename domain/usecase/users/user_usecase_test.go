package users

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tpp/msf/domain/repository/users"
	"github.com/tpp/msf/domain/repository/users/mocks"
	"github.com/tpp/msf/model"
	"github.com/tpp/msf/shared/base"
	"github.com/tpp/msf/shared/utils"
)

func Test_usecase_ListUsers(t *testing.T) {

	baseUsecase := base.NewBaseUsecase("test_user_usecase")
	tests := []struct {
		name   string
		setup  utils.TestSetup[*users.Repository]
		expect func(r *require.Assertions, users []*model.User, total int64, err error)
	}{
		// TODO: Add test cases.
		{
			name: "success",
			setup: func(t *testing.T, _ *require.Assertions, r *users.Repository) utils.TestTeardown {
				repo := mocks.NewRepository(t)
				repo.On("List", mock.Anything, 10, 20, "", base.Filters(nil)).Return(
					[]*model.User{{ID: 1}, {ID: 2}},
					int64(2),
					nil,
				)

				*r = repo
				return func() {}
			},
			expect: func(r *require.Assertions, users []*model.User, total int64, err error) {
				r.EqualValues(2, len(users))
				r.EqualValues([]*model.User{{ID: 1}, {ID: 2}}, users)
				r.EqualValues(2, total)
				r.NoError(err)
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			var repo users.Repository
			var assertion = require.New(t)
			td := tt.setup(t, assertion, &repo)
			if td != nil {
				defer td()
			}
			u := &usecase{
				Usecase:  baseUsecase,
				userRepo: repo,
			}
			users, total, err := u.ListUsers(nil, 10, 20, "", nil)
			tt.expect(assertion, users, total, err)

		})
	}
}
