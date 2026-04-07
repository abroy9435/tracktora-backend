package repository

import (
	"context"
	"encoding/json"
	"tracktora-backend/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ResumeRepository struct {
	DB *pgxpool.Pool
}

func NewResumeRepository(db *pgxpool.Pool) *ResumeRepository {
	return &ResumeRepository{DB: db}
}

// --- VAULT HELPERS ---

func (r *ResumeRepository) AddProject(p models.Project) error {
	_, err := r.DB.Exec(context.Background(), "INSERT INTO projects (user_id, project_name, description, tech_stack, project_url) VALUES ($1, $2, $3, $4, $5)", p.UserID, p.ProjectName, p.Description, p.TechStack, p.ProjectURL)
	return err
}

func (r *ResumeRepository) GetProjectsByUser(userID string) ([]models.Project, error) {
	rows, err := r.DB.Query(context.Background(), "SELECT id, project_name, description, tech_stack, project_url FROM projects WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []models.Project
	for rows.Next() {
		var p models.Project
		rows.Scan(&p.ID, &p.ProjectName, &p.Description, &p.TechStack, &p.ProjectURL)
		res = append(res, p)
	}
	return res, nil
}

func (r *ResumeRepository) AddExperience(e models.Experience) error {
	_, err := r.DB.Exec(context.Background(), "INSERT INTO experiences (user_id, company_name, role_title, start_date, end_date, is_current, bullet_points) VALUES ($1, $2, $3, $4, $5, $6, $7)", e.UserID, e.CompanyName, e.RoleTitle, e.StartDate, e.EndDate, e.IsCurrent, e.BulletPoints)
	return err
}

func (r *ResumeRepository) GetExperiencesByUser(userID string) ([]models.Experience, error) {
	rows, err := r.DB.Query(context.Background(), "SELECT id, company_name, role_title, start_date, end_date, is_current, bullet_points FROM experiences WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []models.Experience
	for rows.Next() {
		var e models.Experience
		rows.Scan(&e.ID, &e.CompanyName, &e.RoleTitle, &e.StartDate, &e.EndDate, &e.IsCurrent, &e.BulletPoints)
		res = append(res, e)
	}
	return res, nil
}

func (r *ResumeRepository) AddEducation(e models.Education) error {
	_, err := r.DB.Exec(context.Background(), "INSERT INTO educations (user_id, institution, degree, field_of_study, start_year, end_year) VALUES ($1, $2, $3, $4, $5, $6)", e.UserID, e.Institution, e.Degree, e.FieldOfStudy, e.StartYear, e.EndYear)
	return err
}

func (r *ResumeRepository) GetEducationsByUser(userID string) ([]models.Education, error) {
	rows, err := r.DB.Query(context.Background(), "SELECT id, institution, degree, field_of_study, start_year, end_year FROM educations WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []models.Education
	for rows.Next() {
		var e models.Education
		rows.Scan(&e.ID, &e.Institution, &e.Degree, &e.FieldOfStudy, &e.StartYear, &e.EndYear)
		res = append(res, e)
	}
	return res, nil
}

func (r *ResumeRepository) AddSkill(s models.Skill) error {
	_, err := r.DB.Exec(context.Background(), "INSERT INTO skills (user_id, skill_name, category) VALUES ($1, $2, $3)", s.UserID, s.SkillName, s.Category)
	return err
}

func (r *ResumeRepository) GetSkillsByUser(userID string) ([]models.Skill, error) {
	rows, err := r.DB.Query(context.Background(), "SELECT id, skill_name, category FROM skills WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []models.Skill
	for rows.Next() {
		var s models.Skill
		rows.Scan(&s.ID, &s.SkillName, &s.Category)
		res = append(res, s)
	}
	return res, nil
}

func (r *ResumeRepository) AddCertification(c models.Certification) error {
	_, err := r.DB.Exec(context.Background(), "INSERT INTO certifications (user_id, name, issuing_organization, issue_date, credential_url) VALUES ($1, $2, $3, $4, $5)", c.UserID, c.Name, c.IssuingOrganization, c.IssueDate, c.CredentialURL)
	return err
}

func (r *ResumeRepository) GetCertificationsByUser(userID string) ([]models.Certification, error) {
	rows, err := r.DB.Query(context.Background(), "SELECT id, name, issuing_organization, issue_date, credential_url FROM certifications WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []models.Certification
	for rows.Next() {
		var c models.Certification
		rows.Scan(&c.ID, &c.Name, &c.IssuingOrganization, &c.IssueDate, &c.CredentialURL)
		res = append(res, c)
	}
	return res, nil
}

// --- VAULT UPDATES & DELETES ---

func (r *ResumeRepository) UpdateProject(p models.Project) error {
	query := `UPDATE projects SET project_name=$1, description=$2, tech_stack=$3, project_url=$4 WHERE id=$5 AND user_id=$6`
	_, err := r.DB.Exec(context.Background(), query, p.ProjectName, p.Description, p.TechStack, p.ProjectURL, p.ID, p.UserID)
	return err
}

func (r *ResumeRepository) DeleteProject(id, userID string) error {
	_, err := r.DB.Exec(context.Background(), "DELETE FROM projects WHERE id=$1 AND user_id=$2", id, userID)
	return err
}

func (r *ResumeRepository) UpdateExperience(e models.Experience) error {
	query := `UPDATE experiences SET company_name=$1, role_title=$2, start_date=$3, end_date=$4, is_current=$5, bullet_points=$6 WHERE id=$7 AND user_id=$8`
	_, err := r.DB.Exec(context.Background(), query, e.CompanyName, e.RoleTitle, e.StartDate, e.EndDate, e.IsCurrent, e.BulletPoints, e.ID, e.UserID)
	return err
}

func (r *ResumeRepository) DeleteExperience(id, userID string) error {
	_, err := r.DB.Exec(context.Background(), "DELETE FROM experiences WHERE id=$1 AND user_id=$2", id, userID)
	return err
}

func (r *ResumeRepository) UpdateEducation(e models.Education) error {
	query := `UPDATE educations SET institution=$1, degree=$2, field_of_study=$3, start_year=$4, end_year=$5 WHERE id=$6 AND user_id=$7`
	_, err := r.DB.Exec(context.Background(), query, e.Institution, e.Degree, e.FieldOfStudy, e.StartYear, e.EndYear, e.ID, e.UserID)
	return err
}

func (r *ResumeRepository) DeleteEducation(id, userID string) error {
	_, err := r.DB.Exec(context.Background(), "DELETE FROM educations WHERE id=$1 AND user_id=$2", id, userID)
	return err
}

func (r *ResumeRepository) UpdateSkill(s models.Skill) error {
	query := `UPDATE skills SET skill_name=$1, category=$2 WHERE id=$3 AND user_id=$4`
	_, err := r.DB.Exec(context.Background(), query, s.SkillName, s.Category, s.ID, s.UserID)
	return err
}

func (r *ResumeRepository) DeleteSkill(id, userID string) error {
	_, err := r.DB.Exec(context.Background(), "DELETE FROM skills WHERE id=$1 AND user_id=$2", id, userID)
	return err
}

func (r *ResumeRepository) UpdateCertification(c models.Certification) error {
	query := `UPDATE certifications SET name=$1, issuing_organization=$2, issue_date=$3, credential_url=$4 WHERE id=$5 AND user_id=$6`
	_, err := r.DB.Exec(context.Background(), query, c.Name, c.IssuingOrganization, c.IssueDate, c.CredentialURL, c.ID, c.UserID)
	return err
}

func (r *ResumeRepository) DeleteCertification(id, userID string) error {
	_, err := r.DB.Exec(context.Background(), "DELETE FROM certifications WHERE id=$1 AND user_id=$2", id, userID)
	return err
}

func (r *ResumeRepository) DeleteResume(id, userID string) error {
	_, err := r.DB.Exec(context.Background(), "DELETE FROM resumes WHERE id=$1 AND user_id=$2", id, userID)
	return err
}

// --- RESUME COMPILATION ---

func (r *ResumeRepository) SaveResume(res models.Resume) error {
	ex, _ := json.Marshal(res.ExperienceIDs)
	pr, _ := json.Marshal(res.ProjectIDs)
	ed, _ := json.Marshal(res.EducationIDs)
	sk, _ := json.Marshal(res.SkillIDs)
	ce, _ := json.Marshal(res.CertificationIDs)
	query := `INSERT INTO resumes (user_id, resume_name, target_role, summary, experience_ids, project_ids, education_ids, skill_ids, certification_ids) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.DB.Exec(context.Background(), query, res.UserID, res.ResumeName, res.TargetRole, res.Summary, ex, pr, ed, sk, ce)
	return err
}

func (r *ResumeRepository) GetResumesByUser(userID string) ([]models.Resume, error) {
	// 1. Updated query to SELECT ALL columns
	query := `SELECT id, resume_name, target_role, summary, experience_ids, project_ids, education_ids, skill_ids, certification_ids 
	          FROM resumes WHERE user_id = $1`

	rows, err := r.DB.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Resume
	for rows.Next() {
		var res models.Resume
		// Temporary holders for the JSONB columns
		var ex, pr, ed, sk, ce []byte

		// 2. Scan all columns into the struct and our temp JSON buffers
		err := rows.Scan(
			&res.ID,
			&res.ResumeName,
			&res.TargetRole,
			&res.Summary,
			&ex, &pr, &ed, &sk, &ce,
		)
		if err != nil {
			return nil, err
		}

		// 3. Unmarshal the JSON bytes back into Go string slices
		json.Unmarshal(ex, &res.ExperienceIDs)
		json.Unmarshal(pr, &res.ProjectIDs)
		json.Unmarshal(ed, &res.EducationIDs)
		json.Unmarshal(sk, &res.SkillIDs)
		json.Unmarshal(ce, &res.CertificationIDs)

		list = append(list, res)
	}
	return list, nil
}

func (r *ResumeRepository) GetCompiledResume(id, userID string) (*models.CompiledResume, error) {
	var c models.CompiledResume
	var ex, pr, ed, sk, ce []byte
	query := `SELECT resume_name, target_role, summary, experience_ids, project_ids, education_ids, skill_ids, certification_ids 
              FROM resumes WHERE id = $1 AND user_id = $2`
	err := r.DB.QueryRow(context.Background(), query, id, userID).Scan(&c.ResumeDetails.ResumeName, &c.ResumeDetails.TargetRole, &c.ResumeDetails.Summary, &ex, &pr, &ed, &sk, &ce)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(ex, &c.ResumeDetails.ExperienceIDs)
	json.Unmarshal(pr, &c.ResumeDetails.ProjectIDs)
	json.Unmarshal(ed, &c.ResumeDetails.EducationIDs)
	json.Unmarshal(sk, &c.ResumeDetails.SkillIDs)
	json.Unmarshal(ce, &c.ResumeDetails.CertificationIDs)

	// HYDRATION: Fetch actual data in user's priority order
	for _, tid := range c.ResumeDetails.ProjectIDs {
		var p models.Project
		if err := r.DB.QueryRow(context.Background(), "SELECT id, project_name, description, tech_stack, project_url FROM projects WHERE id = $1", tid).Scan(&p.ID, &p.ProjectName, &p.Description, &p.TechStack, &p.ProjectURL); err == nil {
			c.Projects = append(c.Projects, p)
		}
	}
	for _, tid := range c.ResumeDetails.ExperienceIDs {
		var e models.Experience
		if err := r.DB.QueryRow(context.Background(), "SELECT id, company_name, role_title, start_date, end_date, is_current, bullet_points FROM experiences WHERE id = $1", tid).Scan(&e.ID, &e.CompanyName, &e.RoleTitle, &e.StartDate, &e.EndDate, &e.IsCurrent, &e.BulletPoints); err == nil {
			c.Experiences = append(c.Experiences, e)
		}
	}
	for _, tid := range c.ResumeDetails.EducationIDs {
		var edu models.Education
		if err := r.DB.QueryRow(context.Background(), "SELECT id, institution, degree, field_of_study, start_year, end_year FROM educations WHERE id = $1", tid).Scan(&edu.ID, &edu.Institution, &edu.Degree, &edu.FieldOfStudy, &edu.StartYear, &edu.EndYear); err == nil {
			c.Educations = append(c.Educations, edu)
		}
	}
	for _, tid := range c.ResumeDetails.SkillIDs {
		var s models.Skill
		if err := r.DB.QueryRow(context.Background(), "SELECT id, skill_name, category FROM skills WHERE id = $1", tid).Scan(&s.ID, &s.SkillName, &s.Category); err == nil {
			c.Skills = append(c.Skills, s)
		}
	}
	for _, tid := range c.ResumeDetails.CertificationIDs {
		var cert models.Certification
		if err := r.DB.QueryRow(context.Background(), "SELECT id, name, issuing_organization, issue_date, credential_url FROM certifications WHERE id = $1", tid).Scan(&cert.ID, &cert.Name, &cert.IssuingOrganization, &cert.IssueDate, &cert.CredentialURL); err == nil {
			c.Certifications = append(c.Certifications, cert)
		}
	}

	return &c, nil
}
