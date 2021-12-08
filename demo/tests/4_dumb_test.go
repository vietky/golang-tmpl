package tests

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/vietky/golang-tmpl/demo/bettercodes"
	"github.com/vietky/golang-tmpl/demo/mocks"
)

func TestAnotherDumb(t *testing.T) {
	ctrl := gomock.NewController(t)
	// defer ctrl.Finish()

	mockLogger := mocks.NewMockILogger(ctrl)
	mockLogger.EXPECT().Log(gomock.Any()).Times(1)

	dumb := bettercodes.NewAnotherDumb(mockLogger)
	dumb.Hello()

	v := dumb.Calculate(1, 2)
	if v != 3 {
		t.Errorf("v should be 3. Current Value is %v", v)
	}
}
