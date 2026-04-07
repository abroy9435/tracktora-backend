package handlers

import (
	"tracktora-backend/internal/database"
	"tracktora-backend/internal/models"
	"tracktora-backend/internal/repository"

	"github.com/gofiber/fiber/v2"
)

func getResumeRepo() *repository.ResumeRepository {
	return repository.NewResumeRepository(database.DB)
}

// --- VAULT: PROJECTS ---
func AddProject(c *fiber.Ctx) error {
	var p models.Project
	if err := c.BodyParser(&p); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	p.UserID = c.Locals("user_id").(string)
	if err := getResumeRepo().AddProject(p); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"message": "Project saved"})
}

func GetProjects(c *fiber.Ctx) error {
	projects, err := getResumeRepo().GetProjectsByUser(c.Locals("user_id").(string))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(projects)
}

func UpdateProject(c *fiber.Ctx) error { return c.SendStatus(501) }
func DeleteProject(c *fiber.Ctx) error { return c.SendStatus(501) }

// --- VAULT: EXPERIENCES ---
func AddExperience(c *fiber.Ctx) error {
	var e models.Experience
	if err := c.BodyParser(&e); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	e.UserID = c.Locals("user_id").(string)
	if err := getResumeRepo().AddExperience(e); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"message": "Experience saved"})
}

func GetExperiences(c *fiber.Ctx) error {
	exps, err := getResumeRepo().GetExperiencesByUser(c.Locals("user_id").(string))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(exps)
}

func UpdateExperience(c *fiber.Ctx) error { return c.SendStatus(501) }
func DeleteExperience(c *fiber.Ctx) error { return c.SendStatus(501) }

// --- VAULT: EDUCATION ---
func AddEducation(c *fiber.Ctx) error {
	var e models.Education
	if err := c.BodyParser(&e); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	e.UserID = c.Locals("user_id").(string)
	if err := getResumeRepo().AddEducation(e); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"message": "Education saved"})
}

func GetEducations(c *fiber.Ctx) error {
	edus, err := getResumeRepo().GetEducationsByUser(c.Locals("user_id").(string))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(edus)
}

func UpdateEducation(c *fiber.Ctx) error { return c.SendStatus(501) }
func DeleteEducation(c *fiber.Ctx) error { return c.SendStatus(501) }

// --- VAULT: SKILLS ---
func AddSkill(c *fiber.Ctx) error {
	var s models.Skill
	if err := c.BodyParser(&s); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	s.UserID = c.Locals("user_id").(string)
	if err := getResumeRepo().AddSkill(s); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"message": "Skill saved"})
}

func GetSkills(c *fiber.Ctx) error {
	skills, err := getResumeRepo().GetSkillsByUser(c.Locals("user_id").(string))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(skills)
}

func UpdateSkill(c *fiber.Ctx) error { return c.SendStatus(501) }
func DeleteSkill(c *fiber.Ctx) error { return c.SendStatus(501) }

// --- VAULT: CERTIFICATIONS (NEW TWEAK) ---
func AddCertification(c *fiber.Ctx) error {
	var cert models.Certification
	if err := c.BodyParser(&cert); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	cert.UserID = c.Locals("user_id").(string)
	if err := getResumeRepo().AddCertification(cert); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"message": "Certification saved"})
}

func GetCertifications(c *fiber.Ctx) error {
	certs, err := getResumeRepo().GetCertificationsByUser(c.Locals("user_id").(string))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(certs)
}

func UpdateCertification(c *fiber.Ctx) error { return c.SendStatus(501) }
func DeleteCertification(c *fiber.Ctx) error { return c.SendStatus(501) }

// --- RESUME MANAGEMENT ---
func SaveResumeBlueprint(c *fiber.Ctx) error {
	var r models.Resume
	if err := c.BodyParser(&r); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	r.UserID = c.Locals("user_id").(string)
	if err := getResumeRepo().SaveResume(r); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"message": "Resume blueprint stored"})
}

func GetSavedResumes(c *fiber.Ctx) error {
	resumes, err := getResumeRepo().GetResumesByUser(c.Locals("user_id").(string))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(resumes)
}

func GetCompiledResume(c *fiber.Ctx) error {
	compiled, err := getResumeRepo().GetCompiledResume(c.Params("id"), c.Locals("user_id").(string))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(compiled)
}

func DeleteResume(c *fiber.Ctx) error { return c.SendStatus(501) }
