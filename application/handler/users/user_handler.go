package users

import (
	"fmt"
	"gorm.io/gorm"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/tpp/msf/model"

	"github.com/go-chi/chi"
	"github.com/tpp/msf/domain/usecase/users"
	"github.com/tpp/msf/shared/base"
)

type Handler interface {
	List(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	CreateRole(w http.ResponseWriter, r *http.Request)
	AssignRole(w http.ResponseWriter, r *http.Request)
	UpdatePassWord(w http.ResponseWriter, r *http.Request)
	UpdateName(w http.ResponseWriter, r *http.Request)
	UpdateActive(w http.ResponseWriter, r *http.Request)
	AdminResetPWForUser(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	base.HTTPHandler
	usecase users.Usecase
}

func (h *handler) List(w http.ResponseWriter, r *http.Request) {

	var reqParams listParams
	ctx, err := h.Parse(r, &reqParams, base.ParseTypeParam)
	if err != nil {
		h.ResponseBadRequest(w, err)
		return
	}
	reqParams.EnsureDefault()
	reqParams.RegisterFilterKeys(listAcceptedFilterKeys)

	if isValidationErrs, err := h.Validate(&reqParams); err != nil {
		if isValidationErrs {
			h.ResponseBadRequest(w, err)
			h.Debug(ctx).Err(err).Msg("InvalidListParams")
			return
		}
		h.ResponseInternalServerError(w)
		h.Error(ctx).Err(err).Msg("InvalidValidatorError")
		return
	}
	h.QueryOnly(ctx)
	user := ctx.User()
	// filters.0.key=full_name&filters.0.operator=like&filters.0.value=Hank
	filters := addFilter(user, reqParams.Filters)
	users, total, err := h.usecase.ListUsers(ctx, reqParams.Limit, reqParams.Offset, reqParams.Sort, filters)
	var listUsers []*model.ListUserRes
	for _, user := range users {
		listUsers = append(listUsers, user.MapUserModel())
	}
	if err != nil {
		h.ResponseInternalServerError(w)
		return
	}
	if total == 0 {
		h.ResponseSuccess(w, "No data to show")
		return
	}
	var totalPages int64
	totalPages = int64(math.Ceil(float64(total) / float64(reqParams.Limit)))
	h.ResponseSuccess(w, &base.Pagination{
		Limit:     reqParams.Limit,
		Offset:    reqParams.Offset,
		Sort:      reqParams.Sort,
		Total:     total,
		TotalPage: totalPages,
		List:      listUsers,
	})
}

func addFilter(user *model.User, filters []*base.Filter) []*base.Filter {
	needFilter := false

	if needFilter {
		return append(filters, &base.Filter{
			Key:      "full_name",
			Operator: "like",
			Value:    fmt.Sprint(user.FullName),
		})
	} else {
		return filters
	}

}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {

	ctx, err := h.Parse(r, nil, base.ParseTypeNone)
	if err != nil {
		h.ResponseBadRequest(w, err)
		return
	}

	var userID uint64
	uStr := chi.URLParamFromCtx(ctx, "id")
	if ctx.User() != nil && fmt.Sprint(ctx.User().ID) == uStr {
		// check if there is same id, do not need to convert from string to uint64
		userID = ctx.User().ID
	} else {
		// other wise, do convertion
		userID, err = strconv.ParseUint(uStr, 10, 64)
		if err != nil {
			h.ResponseBadRequest(w, err)
			h.Debug(ctx).Err(err).Msg("InvalidUserID")
			return
		}
	}

	h.QueryOnly(ctx)
	user, err := h.usecase.GetUser(ctx, userID)
	if err != nil {
		if err == base.ErrorNotFound {
			h.ResponseNotFound(w)
			return
		}
		h.ResponseInternalServerError(w)
		return
	}

	h.ResponseSuccess(w, user.MapUserModel())

}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {

	var createUserReq *createUserReq
	var err error
	ctx, err := h.Parse(r, &createUserReq, base.ParseTypeJSON)
	if err != nil {
		h.ResponseBadRequest(w)
		return
	}

	if isValidationError, err := h.Validate(createUserReq); err != nil {
		if isValidationError {
			h.ResponseBadRequest(w, err)
			h.Debug(ctx).Err(err).Msg("validation error")
			return
		}
		h.ResponseInternalServerError(w)
		h.Error(ctx).Err(err).Msg("InvalidValidationError")
		return
	}

	h.Start(ctx)
	defer func() {
		if err != nil {
			h.Rollback(ctx)
		} else {
			h.Commit(ctx)
		}
	}()

	var user = createUserReq.mapUserModel()
	err = h.usecase.CreateUser(ctx, user)
	if err != nil {
		if _, ok := err.(base.APIErrors); ok {
			h.ResponseBadRequest(w, err)
			return
		}
		if err == base.ErrorNotFound {
			h.ResponseNotFound(w)
			return
		}
		h.ResponseInternalServerError(w)
		return
	}
	var roles = createUserReq.Role
	err = h.usecase.AssignMultipleRole(ctx, int64(user.ID), roles)
	if err != nil {
		h.ResponseInternalServerError(w)
		return
	}

	h.ResponseSuccess(w, "An email will be sent to user's email address with a link to activate the account", 201)

}

func (h *handler) CreateRole(w http.ResponseWriter, r *http.Request) {

	var createRoleReq *createRoleReq
	var err error
	ctx, err := h.Parse(r, &createRoleReq, base.ParseTypeJSON)
	if err != nil {
		h.ResponseBadRequest(w)
		return
	}

	if isValidationError, err := h.Validate(createRoleReq); err != nil {
		if isValidationError {
			h.ResponseBadRequest(w, err)
			h.Debug(ctx).Err(err).Msg("validation error")
			return
		}
		h.ResponseInternalServerError(w)
		h.Error(ctx).Err(err).Msg("InvalidValidationError")
		return
	}

	h.Start(ctx)
	defer func() {
		if err != nil {
			h.Rollback(ctx)
		} else {
			h.Commit(ctx)
		}
	}()

	var role = createRoleReq.mapRoleModel()
	err = h.usecase.CreateRole(ctx, role)
	if err != nil {
		if _, ok := err.(base.APIErrors); ok {
			h.ResponseBadRequest(w, err)
			return
		}
		if err == base.ErrorNotFound {
			h.ResponseNotFound(w)
			return
		}
		h.ResponseInternalServerError(w)
		return
	}

	h.ResponseSuccess(w, "Successfully", 201)

}

func (h *handler) AssignRole(w http.ResponseWriter, r *http.Request) {

	var assignRoleReq *assignRoleReq
	var err error
	ctx, err := h.Parse(r, &assignRoleReq, base.ParseTypeJSON)
	if err != nil {
		h.ResponseBadRequest(w)
		return
	}

	if len(assignRoleReq.RoleIDs) == 0 {
		if checkEmptyString(assignRoleReq.UserName) {
			h.ResponseBadRequest(w, "nothing to change")
			return
		}
	}

	h.Start(ctx)
	defer func() {
		if err != nil {
			h.Rollback(ctx)
		} else {
			h.Commit(ctx)
		}
	}()

	if len(assignRoleReq.RoleIDs) != 0 {
		err = h.usecase.AssignRole(ctx, int64(assignRoleReq.UserID), assignRoleReq.RoleIDs)
		if err != nil {
			if _, ok := err.(base.APIErrors); ok {
				h.ResponseBadRequest(w, err)
				return

			}
			if err == base.ErrorNotFound {
				h.ResponseNotFound(w)
				return
			}
			h.ResponseInternalServerError(w)
			return
		}
	}

	if !checkEmptyString(assignRoleReq.UserName) {
		err = h.usecase.UpdateName(ctx, assignRoleReq.UserName, int64(assignRoleReq.UserID))
		if err != nil {
			if _, ok := err.(base.APIErrors); ok {
				h.ResponseBadRequest(w, err)
				return

			}
			if err == base.ErrorNotFound {
				h.ResponseNotFound(w)
				return
			}
			h.ResponseInternalServerError(w)
			return
		}
	}

	h.ResponseSuccess(w, nil, 204)

}
func (h *handler) UpdatePassWord(w http.ResponseWriter, r *http.Request) {
	var passWord *PassWordReset
	ctx, err := h.Parse(r, &passWord, base.ParseTypeJSON)
	if err != nil {
		h.ResponseBadRequest(w, err)
		return
	}
	if isValidationError, err := h.Validate(passWord); err != nil {
		if isValidationError {
			h.ResponseBadRequest(w, err)
			h.Debug(ctx).Err(err).Msg("InvalidPassword")
			return
		}
		h.ResponseInternalServerError(w)
		h.Error(ctx).Err(err).Msg("InvalidValidatorError")
		return
	}

	h.Start(ctx)
	defer func() {
		if err != nil {
			h.Rollback(ctx)
		} else {
			h.Commit(ctx)
		}
	}()
	err = h.usecase.UpdatePassWord(ctx, int64(ctx.User().ID), passWord.Password)
	if err != nil {
		if err == base.ErrorNotFound {
			h.ResponseNotFound(w)
			return
		}
		h.ResponseInternalServerError(w)
		return
	}

	h.ResponseSuccess(w, "Password is changed")

}
func (h *handler) UpdateName(w http.ResponseWriter, r *http.Request) {
	var name *UpdateName
	ctx, err := h.Parse(r, &name, base.ParseTypeJSON)
	if err != nil {
		h.ResponseBadRequest(w, err)
		return
	}
	var checkName string = name.FullName
	if checkEmptyString(checkName) {
		h.Error(ctx).Msg("Name is required")
		h.ResponseBadRequest(w, err)
		return
	}
	h.QueryOnly(ctx)
	err = h.usecase.UpdateName(ctx, name.FullName, int64(ctx.User().ID))
	if err != nil {
		if err == base.ErrorNotFound {
			h.ResponseNotFound(w)
			return
		}
		h.ResponseInternalServerError(w)
		return
	}

	h.ResponseSuccess(w, "Your name is changed")

}
func (h *handler) UpdateActive(w http.ResponseWriter, r *http.Request) {
	var active *UpdateActive
	ctx, err := h.Parse(r, &active, base.ParseTypeJSON)
	if err != nil {
		h.ResponseBadRequest(w, err)
		return
	}

	h.Start(ctx)
	defer func() {
		if err != nil {
			h.Rollback(ctx)
		} else {
			h.Commit(ctx)
		}
	}()
	err = h.usecase.UpdateActive(ctx, int64(active.UserId), active.IsActive)
	if err != nil {
		if err == base.ErrorNotFound {
			h.ResponseNotFound(w)
			return
		}
		h.ResponseInternalServerError(w)
		return
	}

	h.ResponseSuccess(w, "active is changed")

}
func (h *handler) AdminResetPWForUser(w http.ResponseWriter, r *http.Request) {
	var passWord *PassWordResetByAdmin
	ctx, err := h.Parse(r, &passWord, base.ParseTypeJSON)
	if err != nil {
		h.ResponseBadRequest(w, err)
		return
	}
	var validatePassword *PassWordReset
	validatePassword = validatePassword.setPassWordReset(passWord.Password)
	if isValidationError, err := h.Validate(validatePassword); err != nil {
		if isValidationError {
			h.ResponseBadRequest(w, err)
			h.Debug(ctx).Err(err).Msg("InvalidPassword")
			return
		}
		h.ResponseInternalServerError(w)
		h.Error(ctx).Err(err).Msg("InvalidValidatorError")
		return
	}

	h.Start(ctx)
	defer func() {
		if err != nil {
			h.Rollback(ctx)
		} else {
			h.Commit(ctx)
		}
	}()
	err = h.usecase.UpdatePassWord(ctx, int64(passWord.UserId), passWord.Password)
	if err == gorm.ErrRecordNotFound {
		h.ResponseBadRequest(w, "user is not existed")
	}
	if err != nil {
		if err == base.ErrorNotFound {
			h.ResponseNotFound(w)
			return
		}
		h.ResponseInternalServerError(w)
		return
	}

	h.ResponseSuccess(w, "user's password is changed")

}

func New() Handler {
	return &handler{
		HTTPHandler: base.NewBaseHTTPHandler("users"),
		usecase:     users.New(),
	}
}
func checkEmptyString(str string) bool {
	str = strings.Replace(str, " ", "", -1)
	if str == "" {
		return true
	}
	return false
}
