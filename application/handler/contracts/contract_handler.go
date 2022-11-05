package contracts

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/tpp/msf/domain/usecase/contracts"
	"github.com/tpp/msf/model"
	"github.com/tpp/msf/shared/base"
)

type Handler interface {
	List(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	base.HTTPHandler
	usecase contracts.Usecase
}

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
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
			h.Debug(ctx).Err(err).Msg("IndalidListParams")
			return
		}
		h.ResponseInternalServerError(w)
		h.Error(ctx).Err(err).Msg("InvalidValidatorError")
		return
	}
	h.QueryOnly(ctx)

	user := ctx.User()
	filters := addFilter(user, reqParams.Filters)

	contracts, total, err := h.usecase.List(ctx, reqParams.Limit, reqParams.Offset, reqParams.Sort, filters)
	if err != nil {
		h.ResponseInternalServerError(w)
		h.Error(ctx).Err(err).Msg("InvalidValidationError")
		return
	}
	h.ResponseSuccess(w, &base.Pagination{
		Limit:  reqParams.Limit,
		Offset: reqParams.Offset,
		Total:  total,
		List:   contracts,
	})
}

func addFilter(user *model.User, filters []*base.Filter) []*base.Filter {

	userpermissions := user.Roles[0].Permissions
	needFilter := false
	for _, p := range userpermissions {
		if p.Name == "VIEW_CONTRACT_LIST" {
			needFilter = true
			break
		}
	}

	if needFilter {
		return append(filters, &base.Filter{
			Key:      "supply_vendor_id",
			Operator: "eq",
			Value:    fmt.Sprint(user.OrgID),
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

	var contractID uint64
	uStr := chi.URLParamFromCtx(ctx, "id")

	contractID, err = strconv.ParseUint(uStr, 10, 64)
	if err != nil {
		h.ResponseBadRequest(w, err)
		h.Debug(ctx).Err(err).Msg("InvalidUserID")
		return

	}

	h.QueryOnly(ctx)
	user, err := h.usecase.GetContract(ctx, contractID)
	if err != nil {
		if err == base.ErrorNotFound {
			h.ResponseNotFound(w)
		} else {
			h.ResponseInternalServerError(w)
		}
		return
	}

	h.ResponseSuccess(w, user)

}

func New() Handler {
	return &handler{
		HTTPHandler: base.NewBaseHTTPHandler("contracts"),
		usecase:     contracts.New(),
	}
}
