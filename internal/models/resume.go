package models

import "time"

type Experience struct {
	ID           string     `json:"id"`
	UserID       string     `json:"user_id"`
	CompanyName  string     `json:"company_name"`
	RoleTitle    string     `json:"role_title"`
	StartDate    time.Time  `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	IsCurrent    bool       `json:"is_current"`
	BulletPoints string     `json:"bullet_points"`
}

type Project struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	ProjectName string `json:"project_name"`
	Description string `json:"description"`
	TechStack   string `json:"tech_stack"`
	ProjectURL  string `json:"project_url"`
}

type Education struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	Institution  string `json:"institution"`
	Degree       string `json:"degree"`
	FieldOfStudy string `json:"field_of_study"`
	StartYear    string `json:"start_year"`
	EndYear      string `json:"end_year"`
}

type Skill struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	SkillName string `json:"skill_name"`
	Category  string `json:"category"`
}

type Certification struct {
	ID                  string     `json:"id"`
	UserID              string     `json:"user_id"`
	Name                string     `json:"name"`
	IssuingOrganization string     `json:"issuing_organization"`
	IssueDate           *time.Time `json:"issue_date"`
	CredentialURL       string     `json:"credential_url"`
}

type Resume struct {
	ID               string   `json:"id"`
	UserID           string   `json:"user_id"`
	ResumeName       string   `json:"resume_name"`
	TargetRole       string   `json:"target_role"`
	Summary          string   `json:"summary"` // Tailored summary per resume
	ExperienceIDs    []string `json:"experience_ids"`
	ProjectIDs       []string `json:"project_ids"`
	EducationIDs     []string `json:"education_ids"`
	SkillIDs         []string `json:"skill_ids"`
	CertificationIDs []string `json:"certification_ids"`
}

type CompiledResume struct {
	ResumeDetails  Resume          `json:"resume_details"`
	Experiences    []Experience    `json:"experiences"`
	Projects       []Project       `json:"projects"`
	Educations     []Education     `json:"educations"`
	Skills         []Skill         `json:"skills"`
	Certifications []Certification `json:"certifications"`
}
