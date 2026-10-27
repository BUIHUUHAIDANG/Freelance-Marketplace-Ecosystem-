USE freelance_platform;

-- 1. Freelancer
ALTER TABLE freelancer
  DROP FOREIGN KEY fk_freelancer_user_account,
  ADD CONSTRAINT fk_freelancer_user_account 
    FOREIGN KEY (user_account_id) REFERENCES user_account (id) 
    ON DELETE CASCADE ON UPDATE CASCADE;

-- 2. Hire Manager
ALTER TABLE hire_manager
  DROP FOREIGN KEY fk_hire_manager_user_account,
  DROP FOREIGN KEY fk_hire_manager_company,
  ADD CONSTRAINT fk_hire_manager_user_account 
    FOREIGN KEY (user_account_id) REFERENCES user_account (id) 
    ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT fk_hire_manager_company 
    FOREIGN KEY (company_id) REFERENCES company (id) 
    ON DELETE RESTRICT ON UPDATE CASCADE;

-- 3. Qualification & Test
ALTER TABLE has_skill
  DROP FOREIGN KEY fk_has_skill_freelancer,
  DROP FOREIGN KEY fk_has_skill_skill,
  ADD CONSTRAINT fk_has_skill_freelancer 
    FOREIGN KEY (freelancer_id) REFERENCES freelancer (id) 
    ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT fk_has_skill_skill 
    FOREIGN KEY (skill_id) REFERENCES skill (id) 
    ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE certification
  DROP FOREIGN KEY fk_certification_freelancer,
  ADD CONSTRAINT fk_certification_freelancer 
    FOREIGN KEY (freelancer_id) REFERENCES freelancer (id) 
    ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE test_result
  DROP FOREIGN KEY fk_test_result_freelancer,
  DROP FOREIGN KEY fk_test_result_test,
  ADD CONSTRAINT fk_test_result_freelancer 
    FOREIGN KEY (freelancer_id) REFERENCES freelancer (id) 
    ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT fk_test_result_test 
    FOREIGN KEY (test_id) REFERENCES test (id) 
    ON DELETE CASCADE ON UPDATE CASCADE;

-- 4. Job & Job Required Skill
ALTER TABLE job
  DROP FOREIGN KEY fk_job_hire_manager,
  DROP FOREIGN KEY fk_job_duration,
  DROP FOREIGN KEY fk_job_complexity,
  DROP FOREIGN KEY fk_job_main_skill,
  DROP FOREIGN KEY fk_job_payment_type,
  ADD CONSTRAINT fk_job_hire_manager 
    FOREIGN KEY (hire_manager_id) REFERENCES hire_manager (id) 
    ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT fk_job_duration 
    FOREIGN KEY (duration_id) REFERENCES expected_duration (id) 
    ON DELETE RESTRICT ON UPDATE CASCADE,
  ADD CONSTRAINT fk_job_complexity 
    FOREIGN KEY (complexity_id) REFERENCES complexity (id) 
    ON DELETE RESTRICT ON UPDATE CASCADE,
  ADD CONSTRAINT fk_job_main_skill 
    FOREIGN KEY (main_skill_id) REFERENCES skill (id) 
    ON DELETE RESTRICT ON UPDATE CASCADE,
  ADD CONSTRAINT fk_job_payment_type 
    FOREIGN KEY (payment_type_id) REFERENCES payment_type (id) 
    ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE job_require_skill
  DROP FOREIGN KEY fk_job_require_skill_job,
  DROP FOREIGN KEY fk_job_require_skill_skill,
  ADD CONSTRAINT fk_job_require_skill_job 
    FOREIGN KEY (job_id) REFERENCES job (id) 
    ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT fk_job_require_skill_skill 
    FOREIGN KEY (skill_id) REFERENCES skill (id) 
    ON DELETE CASCADE ON UPDATE CASCADE;

-- 5. Proposal
ALTER TABLE proposal
  DROP FOREIGN KEY fk_proposal_job,
  DROP FOREIGN KEY fk_proposal_freelancer,
  DROP FOREIGN KEY fk_proposal_payment_type,
  DROP FOREIGN KEY fk_proposal_status,
  ADD CONSTRAINT fk_proposal_job 
    FOREIGN KEY (job_id) REFERENCES job (id) 
    ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT fk_proposal_freelancer 
    FOREIGN KEY (freelancer_id) REFERENCES freelancer (id) 
    ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT fk_proposal_payment_type 
    FOREIGN KEY (payment_type_id) REFERENCES payment_type (id) 
    ON DELETE RESTRICT ON UPDATE CASCADE,
  ADD CONSTRAINT fk_proposal_status 
    FOREIGN KEY (current_proposal_status_id) REFERENCES status_catalog (id) 
    ON DELETE RESTRICT ON UPDATE CASCADE;

-- 6. Contract
ALTER TABLE contract
  DROP FOREIGN KEY fk_contract_proposal,
  DROP FOREIGN KEY fk_contract_company,
  DROP FOREIGN KEY fk_contract_freelancer,
  DROP FOREIGN KEY fk_contract_payment_type,
  DROP FOREIGN KEY fk_contract_status,
  ADD CONSTRAINT fk_contract_proposal 
    FOREIGN KEY (proposal_id) REFERENCES proposal (id) 
    ON DELETE RESTRICT ON UPDATE CASCADE,
  ADD CONSTRAINT fk_contract_company 
    FOREIGN KEY (company_id) REFERENCES company (id) 
    ON DELETE RESTRICT ON UPDATE CASCADE,
  ADD CONSTRAINT fk_contract_freelancer 
    FOREIGN KEY (freelancer_id) REFERENCES freelancer (id) 
    ON DELETE RESTRICT ON UPDATE CASCADE,
  ADD CONSTRAINT fk_contract_payment_type 
    FOREIGN KEY (payment_type_id) REFERENCES payment_type (id) 
    ON DELETE RESTRICT ON UPDATE CASCADE,
  ADD CONSTRAINT fk_contract_status 
    FOREIGN KEY (current_contract_status_id) REFERENCES status_catalog (id) 
    ON DELETE RESTRICT ON UPDATE CASCADE;

-- 7. Work Session
ALTER TABLE work_session
  DROP FOREIGN KEY fk_work_session_contract,
  ADD CONSTRAINT fk_work_session_contract 
    FOREIGN KEY (contract_id) REFERENCES contract (id) 
    ON DELETE CASCADE ON UPDATE CASCADE;