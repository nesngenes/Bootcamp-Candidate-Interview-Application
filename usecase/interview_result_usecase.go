package usecase

import (
	"fmt"
	"interview_bootcamp/model"
	"interview_bootcamp/repository"
)

type InterviewResultUseCase interface {
	CreateInterviewResult(payload model.InterviewResult) error
	ListInterviewResult() ([]model.InterviewResult, error)
	GetByIdInterviewResult(id string) (model.InterviewResult, error)
	UpdateInterviewResult(payload model.InterviewResult) error
	DeleteInterviewResult(id string) error
}

type interviewResultUseCase struct {
	repo repository.InterviewResultRepository
}

// membuat hasil interview pertama kali
func (cr *interviewResultUseCase) CreateInterviewResult(payload model.InterviewResult) error {
	//pengecekan kolom tidak boleh kosong
	if payload.Id == "" && payload.InterviewId == "" && payload.ResultId == "" && payload.Note == "" {
		return fmt.Errorf("id, interview_id, result_id, note tidak boleh kosong")
	}

	// pengecekan email tidak boleh sama
	isExistInterviewResult, _ := cr.repo.GetByIdInterviewResult(payload.Id)
	if isExistInterviewResult.Id == payload.Id {
		return fmt.Errorf("id: %s sudah tersedia", payload.Id)
	}

	err := cr.repo.Create(payload)
	if err != nil {
		return fmt.Errorf("failed to create new candidate: %v", err)
	}
	return nil
}

// Menampilkan seluruh hasil interview
func (cr *interviewResultUseCase) ListInterviewResult() ([]model.InterviewResult, error) {
	return cr.repo.ListInterviewResult()
}

// menampilkan hasil interview berdasarkan id
func (cr *interviewResultUseCase) GetByIdInterviewResult(id string) (model.InterviewResult, error) {
	interview_result, err := cr.repo.GetByIdInterviewResult(id)
	if err != nil {
		return model.InterviewResult{}, fmt.Errorf("data dengan id: %s tidak ditemukan", id)
	}
	return interview_result, nil
}

// menghapus data hasil iterview
func (cr *interviewResultUseCase) DeleteInterviewResult(id string) error {
	interview_result, err := cr.GetByIdInterviewResult(id)
	if err != nil {
		return fmt.Errorf("data dengan id: %s tidak ditemukan", id)
	}

	err = cr.repo.DeleteInterviewResult(interview_result.Id)
	if err != nil {
		return fmt.Errorf("gagal menghapus hasil interview: %v", err.Error())
	}
	return nil
}

// melakukan perubahan data pada hasil interview
func (cr *interviewResultUseCase) UpdateInterviewResult(payload model.InterviewResult) error {

	//untuk mengecek apakah kolom interview id sudah diisi
	if payload.InterviewId == "" {
		return fmt.Errorf("kolom interview id harus diisi")
	}

	// pengecekan interview id tidak boleh sama
	isExistInterviewResults, _ := cr.repo.GetByIdInterviewResult(payload.InterviewId)
	if isExistInterviewResults.InterviewId == payload.InterviewId {
		return fmt.Errorf("interview id: %s sudah ada", payload.InterviewId)
	}

	//untuk mengecek apakah data dengan nomor tersebut sudah ada
	isExistInterviewResult, _ := cr.repo.GetByIdInterviewResult(payload.InterviewId)
	if isExistInterviewResult.InterviewId == payload.InterviewId && isExistInterviewResult.Id != payload.Id {
		return fmt.Errorf("data dengan interview id: %s sudah ada", payload.InterviewId)
	}

	//untuk melakukan update pada data dengan nomor sesuai kolom
	err := cr.repo.UpdateInterviewResult(payload)
	if err != nil {
		return fmt.Errorf("gagal memperbarui hasil interview: %v", err)
	}

	return nil
}

func NewInterviewResultUseCase(repo repository.InterviewResultRepository) InterviewResultUseCase {
	return &interviewResultUseCase{repo: repo}
}
