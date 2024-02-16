-- MARKETPLACE SERVICE

-- table to store all marketplace profiles
CREATE TABLE IF NOT EXISTS marketplace_profiles (
  "key" varchar(100) NOT NULL,
  "name" varchar(100) NOT NULL,
  author varchar(100) NOT NULL,
  created_on bigint NOT NULL,
  updated_on bigint NOT NULL DEFAULT 0,
  colour_type int8 NOT NULL DEFAULT 0,
  image_type int8 NOT NULL DEFAULT 0,
  camera_type varchar(32) NOT NULL DEFAULT ''::character varying,
  checksum varchar(1024) NOT NULL DEFAULT ''::character varying,
  model_path varchar(1024) NOT NULL DEFAULT ''::character varying, 
  tags VARCHAR[],
  max_trial_images BIGINT NOT NULL DEFAULT 1000::BIGINT,
  price DOUBLE PRECISION NOT NULL DEFAULT 0::DOUBLE PRECISION,
  account_id VARCHAR(100) NOT NULL DEFAULT ''::VARCHAR,
  commission_percentage DOUBLE PRECISION NOT NULL DEFAULT 0::DOUBLE PRECISION,
  number_of_downloads BIGINT NOT NULL DEFAULT 500::BIGINT,
  bundle_key varchar(100) NOT NULL DEFAULT ''::VARCHAR,
  free_for_max BOOLEAN NOT NULL DEFAULT false::BOOLEAN,
  CONSTRAINT marketplace_profiles_key_key UNIQUE (key)
);

-- table to map users to marketplace profiles
CREATE TABLE IF NOT EXISTS user_marketplace_profiles (
  "user_id" varchar(100) NOT NULL,
  "marketplace_profile_key" varchar(100) NOT NULL,
  deleted BOOLEAN NOT NULL DEFAULT false::BOOLEAN,
  bought BOOLEAN NOT NULL DEFAULT false::BOOLEAN,
  images_edited BIGINT NOT NULL DEFAULT 0::BIGINT,
  downloaded_on BIGINT NOT NULL DEFAULT 0::BIGINT,
  max_trial_images BIGINT NOT NULL DEFAULT 1000::BIGINT,
  CONSTRAINT user_marketplace_profiles_pkey PRIMARY KEY (user_id, marketplace_profile_key),
  CONSTRAINT user_marketplace_profiles_user_id_fkey FOREIGN KEY (marketplace_profile_key) REFERENCES marketplace_profiles (key)
);

-- table to store marketplace profile bundles
CREATE TABLE IF NOT EXISTS marketplace_bundles
(
    key varchar(100) NOT NULL,
    name character varying(100) NOT NULL,
    price double precision NOT NULL DEFAULT 0,
    discount double precision NOT NULL DEFAULT 0,
    tags character varying[],
    status character varying(100) NOT NULL DEFAULT 'active'::character varying,
    created_on bigint NOT NULL,
    CONSTRAINT bundles_pkey PRIMARY KEY (key)
);



-- REFERRAL SERVICE

CREATE TABLE "email_invitation" (
    "id" varchar(100) NOT NULL PRIMARY KEY,
    "email" varchar(100) NOT NULL,
    "sent_on" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "referrer_id" varchar(100) NOT NULL,
    FOREIGN KEY (referrer_id) REFERENCES referrer (user_id)
);

CREATE TABLE "referrer" (
    "user_id" varchar(100) UNIQUE NOT NULL,
    "referral_code" varchar(100) UNIQUE NOT NULL,
    "level" varchar(100) NOT NULL,
    "points" float NOT NULL DEFAULT 0,
    "percent" float NOT NULL DEFAULT 0,
    PRIMARY KEY ("user_id", "referral_code")
);

CREATE TABLE "referrals" (
    "referral_code" varchar(100) NOT NULL,
    "referee_id" varchar(100) NOT NULL,
    "name" varchar(100) NOT NULL,
    "email" varchar(100) UNIQUE NOT NULL,
    "opened_on" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "state" varchar(100) NOT NULL DEFAULT '',
    "active_plan" varchar(100) NOT NULL DEFAULT '',
    "device_id" varchar(100) NOT NULL DEFAULT '',
    "paid_amount" float NOT NULL DEFAULT 0,
    "plan_amount" float NOT NULL DEFAULT 0,
    PRIMARY KEY ("referral_code", "device_id"),
    FOREIGN KEY ("referral_code") REFERENCES "referrer"("referral_code")
);

CREATE TABLE "redeem_points" (
    "transaction_id" varchar(100) PRIMARY KEY,
    "referrer_id" varchar(100) NOT NULL,
    "date" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "source" varchar(100) NOT NULL,
    "balance_before" float NOT NULL,
    "balance_after" float NOT NULL,
    "status" varchar(100) NOT NULL, -- credit or debit
    "fingerprint" varchar(100) NOT NULL DEFAULT '',
    FOREIGN KEY ("referrer_id") REFERENCES "referrer"("user_id")
);

CREATE TABLE "rewards" (
    "transaction_id" varchar(100) PRIMARY KEY,
    "referrer_id" varchar(100) NOT NULL,
    "level" varchar(100) NOT NULL,
    "has_claimed" BOOLEAN DEFAULT false,
    "issued_on" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY ("referrer_id") REFERENCES "referrer"("user_id")
);



-- PROFILES_MANAGER SERVICE

-- table to store user's custom profiles
CREATE TABLE IF NOT EXISTS profiles (
  "key" varchar(100) NOT NULL, -- res_where, res_field
  user_id varchar(100) NOT NULL,
  user_email varchar(100) NOT NULL,
  "name" varchar(100) NOT NULL,
  created_on bigint NOT NULL,
  updated_on bigint NOT NULL DEFAULT 0,
  status varchar(100) NOT NULL DEFAULT 'uploading'::character varying,
  colour_type int8 NOT NULL DEFAULT 0,
  image_type int8 NOT NULL DEFAULT 0,
  trained_images int8 NOT NULL DEFAULT 0,
  staged_images int8 NOT NULL DEFAULT 0,
  current_folder varchar(100) NOT NULL DEFAULT ''::character varying,
  trained_folder varchar(100) NOT NULL DEFAULT ''::character varying,
  prebuilt_poc bool NOT NULL DEFAULT false,
  platform int8 NOT NULL DEFAULT 1,
  deleted bool NOT NULL DEFAULT false,
  cam varchar(100) [],
  CONSTRAINT profiles_key_key UNIQUE (key)
);

-- table to store folders inside a profile [parent: profiles, child: catalogs|bubble_folders]
CREATE TABLE IF NOT EXISTS folders (
  "key" varchar(100) NOT NULL,
  profile_key varchar(100) NOT NULL,
  status varchar(256) NOT NULL DEFAULT 'uploading'::character varying,
  uploaded_on bigint NOT NULL DEFAULT 0,
  model_path varchar(1024) NOT NULL DEFAULT ''::character varying,
  "current_catalog" varchar(100) NOT NULL DEFAULT ''::character varying,
  number_of_images bigint NOT NULL DEFAULT 0,
  app_version varchar(32) NOT NULL DEFAULT ''::character varying,
  checksum varchar(1024) NOT NULL DEFAULT ''::character varying,
  cam varchar(100) [],
  CONSTRAINT folders_key_key UNIQUE (key),
  CONSTRAINT folders_profile FOREIGN KEY ("profile_key") REFERENCES "profiles" ("key") ON DELETE CASCADE
);

-- table to store catalogs inside a folder [parent: folders]
CREATE TABLE IF NOT EXISTS catalogs (
  "key" varchar(100) NOT NULL,
  folder_key varchar(100) NOT NULL,
  zip_count int8 NOT NULL,
  "offset" int8 NOT NULL DEFAULT 0,
  status varchar(32) NOT NULL DEFAULT 'not-started'::character varying,
  "name" varchar(512) NOT NULL DEFAULT ''::character varying,
  CONSTRAINT catalogs_key_key UNIQUE (key),
  CONSTRAINT catalogs_folder FOREIGN KEY ("folder_key") REFERENCES "folders" ("key") ON DELETE CASCADE
);

-- table to store the profile training queue
CREATE TABLE IF NOT EXISTS training_queue (
  folder_key varchar(100) NOT NULL,
  queue varchar(256) NOT NULL
);

-- table to store the color profiles for each profile [parent: profiles]
CREATE TABLE IF NOT EXISTS bubbles_profile (
  key varchar(100) NOT NULL,
  profile_key varchar(100) NOT NULL,
  staged_images bigint NOT NULL DEFAULT 0,
  trained_images bigint NOT NULL DEFAULT 0,
  user_id varchar(100) NOT NULL DEFAULT '',
  PRIMARY KEY (key, profile_key)
);

CREATE TABLE IF NOT EXISTS bubbles_folder (
  key varchar(100) NOT NULL,
  folder_key varchar(100) NOT NULL,
  num bigint NOT NULL DEFAULT 0,
  PRIMARY KEY (key, folder_key),
  FOREIGN KEY (folder_key) REFERENCES folders (key) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_training_queue (
    queue varchar(100) NOT NULL,
    user_id varchar(100) PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS folder_training_queue (
    folder_key varchar(100) PRIMARY KEY,
    queue varchar(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS profile_training_queue (
    profile_key varchar(100) PRIMARY KEY,
    queue varchar(100) NOT NULL
);


