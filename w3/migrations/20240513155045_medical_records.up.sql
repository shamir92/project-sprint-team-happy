CREATE TABLE IF NOT EXISTS public.medical_records(
    id serial NOT NULL,
    patient_id character(16) NOT NULL,
    symtomps character varying(2000) NOT NULL,
    medications character varying(2000) NOT NULL,
    created_by uuid NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT medical_records_pkey PRIMARY KEY (id),
    CONSTRAINT medical_records_patient_id_fk 
        FOREIGN KEY (patient_id) 
        REFERENCES public.patients(id)
        ON DELETE CASCADE,
    CONSTRAINT medical_records_created_by_fk 
        FOREIGN KEY (created_by) 
        REFERENCES public.users(id)
        ON DELETE CASCADE
);

CREATE INDEX medical_records_patient_id_idx ON medical_records (patient_id);
CREATE INDEX medical_records_created_by_idx ON medical_records (created_by);