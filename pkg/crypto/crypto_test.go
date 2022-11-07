package crypto

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHashPassword(t *testing.T) {
	Convey("HashPassword", t, func() {
		rawPassword := "sample"
		hashedPassword, err := HashPassword(rawPassword)
		So(err, ShouldBeNil)

		Convey("CheckHashPassword", func() {
			check := CheckPasswordHash(rawPassword, hashedPassword)
			So(check, ShouldBeTrue)
		})
	})
}
