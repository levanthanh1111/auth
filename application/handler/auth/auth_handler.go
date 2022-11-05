package auth

import (
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"

	"github.com/tpp/msf/domain/usecase/auth"
	"github.com/tpp/msf/model"
	"github.com/tpp/msf/shared/base"
)

type Handler interface {
	ListRoles(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	ForgotPassword(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	ListPermissions(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	base.HTTPHandler
	usecase auth.Usecase
}

func (h *handler) ListRoles(w http.ResponseWriter, r *http.Request) {
	var reqParams base.ListParams
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
	//role := ctx.Role()
	// filters.0.key=full_name&filters.0.operator=like&filters.0.value=Hank
	roles, err := h.usecase.ListRoles(ctx, reqParams.Limit, reqParams.Offset, reqParams.Sort, reqParams.Filters)
	var listRoles []*model.ListRoleRes
	for _, role := range roles {
		listRoles = append(listRoles, role.MapRoleListModel())
	}
	if err != nil {
		h.ResponseInternalServerError(w)
		return
	}
	var totalPages int64
	var total int64
	h.ResponseSuccess(w, &base.Pagination{
		Limit:     reqParams.Limit,
		Offset:    reqParams.Offset,
		Sort:      reqParams.Sort,
		Total:     total,
		TotalPage: totalPages,
		List:      listRoles,
	})

}

func (h *handler) ListPermissions(w http.ResponseWriter, r *http.Request) {
	var reqParams base.ListParams
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
	permissions, err := h.usecase.ListPermission(ctx)
	var listPermission []*model.ListPermission
	for _, permission := range permissions {
		var permissonNew = permission.MapPermissionListModel()
		listPermission = append(listPermission, &permissonNew)
	}
	if err != nil {
		h.ResponseInternalServerError(w)
		return
	}
	h.ResponseSuccess(w, listPermission)
}

func addFilter(role *model.Role, filters []*base.Filter) []*base.Filter {
	needFilter := false

	if needFilter {
		return append(filters, &base.Filter{
			Key:      "name",
			Operator: "like",
			Value:    fmt.Sprint(role.Name),
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
	var roleID uint64
	uStr := chi.URLParamFromCtx(ctx, "id")
	if ctx.User() != nil && fmt.Sprint(ctx.User().ID) == uStr {
		roleID = ctx.User().ID
	} else {
		roleID, err = strconv.ParseUint(uStr, 10, 64)
		if err != nil {
			h.ResponseBadRequest(w, err)
			h.Debug(ctx).Err(err).Msg("InvalidRoleID")
			return
		}
	}

	h.QueryOnly(ctx)
	role, err := h.usecase.GetRole(ctx, roleID)
	if err != nil {
		if err == base.ErrorNotFound {
			h.ResponseNotFound(w)
			return
		}
		h.ResponseInternalServerError(w)
		return
	}

	h.ResponseSuccess(w, role.MapRoleListModel())

}
func (h *handler) Login(w http.ResponseWriter, r *http.Request) {

	var loginReq *loginReq
	ctx, err := h.Parse(r, &loginReq, base.ParseTypeJSON)
	if err != nil {
		h.ResponseBadRequest(w, err)
		return
	}

	if isValidationError, err := h.Validate(loginReq); err != nil {
		if isValidationError {
			h.ResponseBadRequest(w, err)
			h.Debug(ctx).Err(err).Msg("InvalidLoginReq")
			return
		}
		h.ResponseInternalServerError(w)
		h.Error(ctx).Err(err).Msg("InvalidValidatorError")
		return
	}

	h.QueryOnly(ctx)
	user, accessToken, err := h.usecase.Login(ctx, loginReq.Email, loginReq.Password)
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

	h.ResponseSuccess(w, userRes{AccessToken: accessToken, Data: user})

}

// forgot password
func (h *handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {

	var sendMailReq *sendMailReq
	ctx, err := h.Parse(r, &sendMailReq, base.ParseTypeJSON)
	if err != nil {
		h.ResponseBadRequest(w, err)
		return
	}

	if isValidationError, err := h.Validate(sendMailReq); err != nil {
		if isValidationError {
			h.ResponseBadRequest(w, err)
			h.Debug(ctx).Err(err).Msg("InvalidLoginReq")
			return
		}
		h.ResponseInternalServerError(w)
		h.Error(ctx).Err(err).Msg("InvalidValidatorError")
		return
	}

	h.QueryOnly(ctx)
	message, err := h.usecase.ForgotPassword(ctx, sendMailReq.Email)
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

	h.ResponseSuccess(w, emailRes{Message: message})

}
func New() Handler {
	return &handler{
		HTTPHandler: base.NewBaseHTTPHandler("auth"),
		usecase:     auth.New(),
	}
}
