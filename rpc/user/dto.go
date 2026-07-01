package rpcuser

import (
	"time"

	"github.com/dingzijian9527-del/Travel-Assistant/kitex_gen/base"
	user "github.com/dingzijian9527-del/Travel-Assistant/kitex_gen/user"
)

func toUserDTO(input *userModel) *user.UserInfo {
	if input == nil {
		return nil
	}
	return &user.UserInfo{
		Id:            input.ID,
		Phone:         input.Phone,
		Nickname:      input.Nickname,
		AvatarUrl:     optionalString(input.AvatarURL),
		HomeCity:      optionalString(input.HomeCity),
		CurrentCity:   optionalString(input.CurrentCity),
		MemberLevel:   input.MemberLevel,
		AccountStatus: input.AccountStatus,
		CreatedAt:     optionalTime(input.CreatedAt),
		UpdatedAt:     optionalTime(input.UpdatedAt),
		DeletedAt:     optionalTimePtr(input.DeletedAt),
	}
}

func successBaseResp() *base.BaseResp {
	return &base.BaseResp{Code: int32(base.ErrorCode_SUCCESS), Msg: "success"}
}

func errorBaseResp(err *serviceError) *base.BaseResp {
	if err == nil {
		return successBaseResp()
	}
	return &base.BaseResp{Code: int32(err.code), Msg: err.message}
}

func registerErrorResponse(err *serviceError) *user.RegisterResponse {
	return &user.RegisterResponse{BaseResp: errorBaseResp(err)}
}

func loginErrorResponse(err *serviceError) *user.LoginResponse {
	return &user.LoginResponse{BaseResp: errorBaseResp(err)}
}

func profileErrorResponse(err *serviceError) *user.GetProfileResponse {
	return &user.GetProfileResponse{BaseResp: errorBaseResp(err)}
}

func updateProfileErrorResponse(err *serviceError) *user.UpdateProfileResponse {
	return &user.UpdateProfileResponse{BaseResp: errorBaseResp(err)}
}

func optionalString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func optionalTime(value time.Time) *string {
	if value.IsZero() {
		return nil
	}
	formatted := value.Format("2006-01-02 15:04:05")
	return &formatted
}

func optionalTimePtr(value *time.Time) *string {
	if value == nil {
		return nil
	}
	return optionalTime(*value)
}
