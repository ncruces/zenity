package zenity

import "testing"

func TestError(t *testing.T) {
	res, err := Error("An error has occured.", Title("Error"), Icon(ErrorIcon))

	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%#v", res)
	}
}

func TestInfo(t *testing.T) {
	res, err := Info("All updates are complete.", Title("Information"), Icon(InfoIcon))

	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%#v", res)
	}
}

func TestWarning(t *testing.T) {
	res, err := Warning("Are you sure you want to proceed?", Title("Warning"), Icon(WarningIcon))

	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%#v", res)
	}
}

func TestQuestion(t *testing.T) {
	res, err := Question("Are you sure you want to proceed?", Title("Question"), Icon(QuestionIcon))

	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%#v", res)
	}
}
