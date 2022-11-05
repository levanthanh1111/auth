//go:generate mockery --name=Repository
package users

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tpp/msf/model"
	"github.com/tpp/msf/shared/base"
	"github.com/tpp/msf/shared/context"
	"github.com/tpp/msf/shared/utils"
	"github.com/stretchr/testify/require"
)

var usersColumns = []string{
	"id",
	"full_name",
	"email",
	"org_id",
}

var orgsColumns = []string{
	"id",
	"name",
	"type",
}

var rolesColumns = []string{
	"id",
	"name",
}

func listUserSuccessMock(mock sqlmock.Sqlmock) {
	mock.ExpectPrepare(`SELECT count(*) FROM "users" LIMIT 10`).ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
	mock.ExpectPrepare(`SELECT * FROM "users" LIMIT 10`).ExpectQuery().WillReturnRows(
		sqlmock.NewRows(usersColumns).
			AddRow(1, "test full name 1", "test.email1@email.test", 1).
			AddRow(2, "test full name 2", "test.email2@email.test", 1),
	)
	mock.ExpectPrepare(`SELECT * FROM "orgs" WHERE "orgs"."id" = $1`).ExpectQuery().WillReturnRows(
		sqlmock.NewRows(orgsColumns).AddRow(1, "test org name", 1),
	)
	mock.ExpectPrepare(`SELECT * FROM "user_role" WHERE "user_role"."user_id" IN ($1,$2)`).ExpectQuery().WithArgs(1, 2).WillReturnRows(
		sqlmock.NewRows(rolesColumns).AddRow(1, 1),
		sqlmock.NewRows(rolesColumns).AddRow(2, 1),
	)
}

func Test_repo_List(t *testing.T) {

	baseRepo := base.NewBaseRepository("test_repository")

	type args struct {
		limit   int
		offset  int
		filters base.Filters
	}
	tests := []struct {
		name   string
		setup  utils.TestSetup[*context.Context] // setup -> teardown
		args   args
		want   []*model.User
		expect func(t *require.Assertions, users []*model.User, total int64, err error)
	}{
		// TODO: Add test cases.
		{
			name: "get success",
			setup: func(t *testing.T, r *require.Assertions, testCtx *context.Context) utils.TestTeardown {

				ctx, mock, cncl, err := utils.NewContextWithDB()
				r.NoErrorf(err, "expected new context with db success")

				listUserSuccessMock(mock)

				*testCtx = ctx

				return cncl
			},
			args: args{
				limit:  10,
				offset: 0,
			},
			expect: func(t *require.Assertions, users []*model.User, total int64, err error) {
				t.EqualValues(2, len(users))
				t.EqualValues(1, users[0].ID, 1)
				t.EqualValues(2, users[1].ID, 2)
				t.EqualValues(1, users[0].OrgID)
				t.NotEmpty(users[0].Org)
				t.NotEmpty(users[1].Org)
				t.EqualValues(2, total)
				t.NoError(err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repo{
				Repository: baseRepo,
			}

			assertion := require.New(t)
			var ctx context.Context
			td := tt.setup(t, assertion, &ctx)
			if td != nil {
				defer td()
			}

			users, total, err := r.List(ctx, tt.args.limit, tt.args.offset, tt.args.filters)

			tt.expect(assertion, users, total, err)
		})
	}
}
