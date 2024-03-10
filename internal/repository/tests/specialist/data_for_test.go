package test_specialist

import (
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/models"
	"github.com/guregu/null"
)

var (
	testcaseSpecialistCreate = []models.SpecialistCreate{
		{
			SpecialistBase: models.SpecialistBase{
				Login:    "testlogin1",
				Password: "testpassword1",
				Fullname: null.NewString("Олег Брахмагуптович", true),
			},
			PhotoUrl: null.NewString("testurl1", true),
		},
		{
			SpecialistBase: models.SpecialistBase{
				Login:    "testlogin2",
				Password: "testpassword2",
				Fullname: null.NewString("", false), // Fullname отсутствует
			},
			PhotoUrl: null.NewString("testurl2", true),
		},
		{
			SpecialistBase: models.SpecialistBase{
				Login:    "testlogin3",
				Password: "testpassword3",
				Fullname: null.NewString("Иван Иванович", true),
			},
			PhotoUrl: null.NewString("", false), // PhotoUrl отсутствует
		},
		{
			SpecialistBase: models.SpecialistBase{
				Login:    "testlogin4",
				Password: "testpassword4",
				Fullname: null.NewString("", false), // Fullname отсутствует
			},
			PhotoUrl: null.NewString("", false), // PhotoUrl отсутствует
		},
	}
	testcaseSpecialistUpdate = []models.Specialist{
		{
			SpecialistCreate: models.SpecialistCreate{
				SpecialistBase: models.SpecialistBase{
					Login:    "testlogin1",
					Password: "testpassword1",
					Fullname: null.NewString("Олег Брахмагуптович", true),
				},
				PhotoUrl: null.NewString("testurl1", true),
			},
			ID:         -1,
			Level:      1,
			IsVerified: true,
		},

		{
			SpecialistCreate: models.SpecialistCreate{
				SpecialistBase: models.SpecialistBase{
					Login:    "testlogin2",
					Password: "testpassword2",
					Fullname: null.NewString("", false), // Fullname отсутствует
				},
				PhotoUrl: null.NewString("testurl2", true),
			},
			ID:         -1,
			Level:      2,
			IsVerified: true,
		},
		{

			SpecialistCreate: models.SpecialistCreate{
				SpecialistBase: models.SpecialistBase{
					Login:    "testlogin3",
					Password: "testpassword3",
					Fullname: null.NewString("Иван Иванович", true),
				},
				PhotoUrl: null.NewString("", false), // PhotoUrl отсутствует
			},
			ID:         -1,
			Level:      1,    // Примерное значение, адаптируйте под ваши нужды
			IsVerified: true, // Примерное значение, адаптируйте под ваши нужды
		},

		{

			SpecialistCreate: models.SpecialistCreate{
				SpecialistBase: models.SpecialistBase{
					Login:    "testlogin4",
					Password: "testpassword4",
					Fullname: null.NewString("", false), // Fullname отсутствует
				},
				PhotoUrl: null.NewString("", false), // PhotoUrl отсутствует
			},
			ID:         -1,
			Level:      0,
			IsVerified: false,
		},
	}
)
