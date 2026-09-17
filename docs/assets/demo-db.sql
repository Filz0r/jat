--
-- PostgreSQL database dump
--

\restrict oEkl51gnaVYxkQZT9dBB3cDhbdVVjtb2OZSNKmn2pbLG93b1AtuKoVrXuM0avgE

-- Dumped from database version 18.3 (Debian 18.3-1.pgdg13+1)
-- Dumped by pg_dump version 18.3 (Debian 18.3-1.pgdg13+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: application_notes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.application_notes (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    body text,
    application_id bigint,
    status_id bigint,
    user_id uuid
);


--
-- Name: application_notes_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.application_notes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: application_notes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.application_notes_id_seq OWNED BY public.application_notes.id;


--
-- Name: application_statuses; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.application_statuses (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    status text,
    user_id uuid,
    kind text NOT NULL,
    archived boolean DEFAULT false NOT NULL
);


--
-- Name: application_statuses_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.application_statuses_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: application_statuses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.application_statuses_id_seq OWNED BY public.application_statuses.id;


--
-- Name: ban_lists; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ban_lists (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    user_id uuid NOT NULL,
    banned_at timestamp with time zone NOT NULL,
    ban_duration bigint NOT NULL,
    ban_reason text NOT NULL,
    active boolean DEFAULT true,
    banned_by uuid NOT NULL
);


--
-- Name: ban_lists_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.ban_lists_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: ban_lists_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.ban_lists_id_seq OWNED BY public.ban_lists.id;


--
-- Name: companies; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.companies (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name text,
    website text,
    created_by uuid NOT NULL,
    edited_by uuid NOT NULL
);


--
-- Name: companies_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.companies_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: companies_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.companies_id_seq OWNED BY public.companies.id;


--
-- Name: company_change_histories; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.company_change_histories (
    id bigint NOT NULL,
    company_id bigint NOT NULL,
    old_name_value text,
    new_name_value text,
    old_website_value text,
    new_website_value text,
    changed_by uuid NOT NULL,
    created_at timestamp with time zone,
    reverted boolean DEFAULT false
);


--
-- Name: company_change_histories_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.company_change_histories_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: company_change_histories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.company_change_histories_id_seq OWNED BY public.company_change_histories.id;


--
-- Name: configs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.configs (
    key text NOT NULL,
    value text NOT NULL,
    type text NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


--
-- Name: job_applications; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.job_applications (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    title text,
    url text,
    user_id uuid,
    status_id bigint,
    company_id bigint
);


--
-- Name: job_applications_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.job_applications_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: job_applications_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.job_applications_id_seq OWNED BY public.job_applications.id;


--
-- Name: refresh_tokens; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.refresh_tokens (
    token text NOT NULL,
    user_id uuid,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    expires_at timestamp with time zone,
    expired_at timestamp with time zone
);


--
-- Name: status_histories; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.status_histories (
    id bigint NOT NULL,
    application_id bigint,
    new_status_id bigint,
    old_status_id bigint,
    created_at timestamp with time zone
);


--
-- Name: status_histories_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.status_histories_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: status_histories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.status_histories_id_seq OWNED BY public.status_histories.id;


--
-- Name: user_settings; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_settings (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    user_id uuid NOT NULL,
    setup_step bigint DEFAULT 1,
    default_application_status_id bigint,
    is_banned boolean DEFAULT false,
    is_enabled boolean DEFAULT false,
    is_admin boolean DEFAULT false
);


--
-- Name: user_settings_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.user_settings_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: user_settings_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.user_settings_id_seq OWNED BY public.user_settings.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    email text,
    username text,
    password text,
    user_settings_id bigint
);


--
-- Name: application_notes id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.application_notes ALTER COLUMN id SET DEFAULT nextval('public.application_notes_id_seq'::regclass);


--
-- Name: application_statuses id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.application_statuses ALTER COLUMN id SET DEFAULT nextval('public.application_statuses_id_seq'::regclass);


--
-- Name: ban_lists id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ban_lists ALTER COLUMN id SET DEFAULT nextval('public.ban_lists_id_seq'::regclass);


--
-- Name: companies id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.companies ALTER COLUMN id SET DEFAULT nextval('public.companies_id_seq'::regclass);


--
-- Name: company_change_histories id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.company_change_histories ALTER COLUMN id SET DEFAULT nextval('public.company_change_histories_id_seq'::regclass);


--
-- Name: job_applications id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.job_applications ALTER COLUMN id SET DEFAULT nextval('public.job_applications_id_seq'::regclass);


--
-- Name: status_histories id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.status_histories ALTER COLUMN id SET DEFAULT nextval('public.status_histories_id_seq'::regclass);


--
-- Name: user_settings id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_settings ALTER COLUMN id SET DEFAULT nextval('public.user_settings_id_seq'::regclass);


--
-- Data for Name: application_notes; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.application_notes (id, created_at, updated_at, deleted_at, body, application_id, status_id, user_id) FROM stdin;
1	2026-09-17 11:28:09.123506+00	2026-09-17 11:28:09.123506+00	\N	Example note 1	2	11	f829a7ff-9d6b-46f8-85e2-95680bde51e3
2	2026-09-17 11:28:37.027984+00	2026-09-17 11:28:37.027984+00	\N	Contacted by recruiter for screening call\n	2	14	f829a7ff-9d6b-46f8-85e2-95680bde51e3
3	2026-09-17 11:29:00.2683+00	2026-09-17 11:29:00.2683+00	\N	Never heard back	2	12	f829a7ff-9d6b-46f8-85e2-95680bde51e3
4	2026-09-17 11:29:26.424878+00	2026-09-17 11:29:26.424878+00	\N	Contacted for screening call	3	14	f829a7ff-9d6b-46f8-85e2-95680bde51e3
5	2026-09-17 11:30:04.670829+00	2026-09-17 11:30:04.670829+00	\N	Rejected by an automatic email after 2 hours	1	13	f829a7ff-9d6b-46f8-85e2-95680bde51e3
\.


--
-- Data for Name: application_statuses; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.application_statuses (id, created_at, updated_at, deleted_at, status, user_id, kind, archived) FROM stdin;
1	2026-09-17 11:25:39.728039+00	2026-09-17 11:25:39.728039+00	\N	Didn't Apply	eba72c65-3e54-444a-8bad-d1ce58fe270f	irrelevant	f
2	2026-09-17 11:25:39.728039+00	2026-09-17 11:25:39.728039+00	\N	Applied	eba72c65-3e54-444a-8bad-d1ce58fe270f	applied	f
3	2026-09-17 11:25:39.728039+00	2026-09-17 11:25:39.728039+00	\N	Ghosted	eba72c65-3e54-444a-8bad-d1ce58fe270f	ghosted	f
4	2026-09-17 11:25:39.728039+00	2026-09-17 11:25:39.728039+00	\N	Rejected	eba72c65-3e54-444a-8bad-d1ce58fe270f	rejected	f
5	2026-09-17 11:25:39.728039+00	2026-09-17 11:25:39.728039+00	\N	Contacted	eba72c65-3e54-444a-8bad-d1ce58fe270f	interviewed	f
6	2026-09-17 11:25:39.728039+00	2026-09-17 11:25:39.728039+00	\N	1st Interview	eba72c65-3e54-444a-8bad-d1ce58fe270f	interviewed	f
7	2026-09-17 11:25:39.728039+00	2026-09-17 11:25:39.728039+00	\N	2nd Interview	eba72c65-3e54-444a-8bad-d1ce58fe270f	interviewed	f
8	2026-09-17 11:25:39.728039+00	2026-09-17 11:25:39.728039+00	\N	3rd Interview	eba72c65-3e54-444a-8bad-d1ce58fe270f	interviewed	f
9	2026-09-17 11:25:39.728039+00	2026-09-17 11:25:39.728039+00	\N	Accepted	eba72c65-3e54-444a-8bad-d1ce58fe270f	accepted	f
10	2026-09-17 11:26:39.304823+00	2026-09-17 11:26:39.304823+00	\N	Didn't Apply	f829a7ff-9d6b-46f8-85e2-95680bde51e3	irrelevant	f
11	2026-09-17 11:26:39.304823+00	2026-09-17 11:26:39.304823+00	\N	Applied	f829a7ff-9d6b-46f8-85e2-95680bde51e3	applied	f
12	2026-09-17 11:26:39.304823+00	2026-09-17 11:26:39.304823+00	\N	Ghosted	f829a7ff-9d6b-46f8-85e2-95680bde51e3	ghosted	f
13	2026-09-17 11:26:39.304823+00	2026-09-17 11:26:39.304823+00	\N	Rejected	f829a7ff-9d6b-46f8-85e2-95680bde51e3	rejected	f
14	2026-09-17 11:26:39.304823+00	2026-09-17 11:26:39.304823+00	\N	Contacted	f829a7ff-9d6b-46f8-85e2-95680bde51e3	interviewed	f
15	2026-09-17 11:26:39.304823+00	2026-09-17 11:26:39.304823+00	\N	1st Interview	f829a7ff-9d6b-46f8-85e2-95680bde51e3	interviewed	f
16	2026-09-17 11:26:39.304823+00	2026-09-17 11:26:39.304823+00	\N	2nd Interview	f829a7ff-9d6b-46f8-85e2-95680bde51e3	interviewed	f
17	2026-09-17 11:26:39.304823+00	2026-09-17 11:26:39.304823+00	\N	3rd Interview	f829a7ff-9d6b-46f8-85e2-95680bde51e3	interviewed	f
18	2026-09-17 11:26:39.304823+00	2026-09-17 11:26:39.304823+00	\N	Accepted	f829a7ff-9d6b-46f8-85e2-95680bde51e3	accepted	f
\.


--
-- Data for Name: ban_lists; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.ban_lists (id, created_at, updated_at, deleted_at, user_id, banned_at, ban_duration, ban_reason, active, banned_by) FROM stdin;
\.


--
-- Data for Name: companies; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.companies (id, created_at, updated_at, deleted_at, name, website, created_by, edited_by) FROM stdin;
1	2026-09-17 11:25:54.156081+00	2026-09-17 11:25:54.156081+00	\N	Test Company 1	\N	eba72c65-3e54-444a-8bad-d1ce58fe270f	eba72c65-3e54-444a-8bad-d1ce58fe270f
2	2026-09-17 11:25:58.646343+00	2026-09-17 11:25:58.646343+00	\N	Test Company 2	\N	eba72c65-3e54-444a-8bad-d1ce58fe270f	eba72c65-3e54-444a-8bad-d1ce58fe270f
3	2026-09-17 11:26:02.694903+00	2026-09-17 11:26:02.694903+00	\N	Test Company 3	\N	eba72c65-3e54-444a-8bad-d1ce58fe270f	eba72c65-3e54-444a-8bad-d1ce58fe270f
4	2026-09-17 11:26:07.729795+00	2026-09-17 11:26:07.729795+00	\N	Test Company 4	\N	eba72c65-3e54-444a-8bad-d1ce58fe270f	eba72c65-3e54-444a-8bad-d1ce58fe270f
5	2026-09-17 11:26:22.741357+00	2026-09-17 11:26:22.741357+00	\N	Test Company 5	https://company5.example.com	eba72c65-3e54-444a-8bad-d1ce58fe270f	eba72c65-3e54-444a-8bad-d1ce58fe270f
\.


--
-- Data for Name: company_change_histories; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.company_change_histories (id, company_id, old_name_value, new_name_value, old_website_value, new_website_value, changed_by, created_at, reverted) FROM stdin;
\.


--
-- Data for Name: configs; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.configs (key, value, type, created_at, updated_at, deleted_at) FROM stdin;
applicationMode	server	jat_mode	2026-09-17 11:25:35.104318+00	2026-09-17 11:25:35.104318+00	\N
initialVersion	dev	string	2026-09-17 11:25:35.106023+00	2026-09-17 11:25:35.106023+00	\N
currentVersion	dev	string	2026-09-17 11:25:35.106855+00	2026-09-17 11:25:35.106855+00	\N
firstAdmin	eba72c65-3e54-444a-8bad-d1ce58fe270f	uuid	2026-09-17 11:25:35.105094+00	2026-09-17 11:25:39.73252+00	\N
isInitialized	true	boolean	2026-09-17 11:25:35.103312+00	2026-09-17 11:25:42.453475+00	\N
\.


--
-- Data for Name: job_applications; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.job_applications (id, created_at, updated_at, deleted_at, title, url, user_id, status_id, company_id) FROM stdin;
4	2026-09-17 11:27:46.16+00	2026-09-17 11:27:54.724305+00	\N	Junior Software Developer	https://example4.com	f829a7ff-9d6b-46f8-85e2-95680bde51e3	11	4
2	2026-09-17 11:27:09.116+00	2026-09-17 11:28:48.117266+00	\N	Junior Software Developer	https://example2.com	f829a7ff-9d6b-46f8-85e2-95680bde51e3	12	2
3	2026-09-17 11:27:29.675+00	2026-09-17 11:29:17.597419+00	\N	Junior Software Developer	https://example3.com	f829a7ff-9d6b-46f8-85e2-95680bde51e3	14	3
1	2026-09-17 11:26:43.88+00	2026-09-17 11:29:44.270139+00	\N	Junior Software Developer	https://example1.com	f829a7ff-9d6b-46f8-85e2-95680bde51e3	13	1
\.


--
-- Data for Name: refresh_tokens; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.refresh_tokens (token, user_id, created_at, updated_at, expires_at, expired_at) FROM stdin;
1fb9c6b94b7ece3e74c6ba01d94c339db3a006786a642681d1101d90400f0550	eba72c65-3e54-444a-8bad-d1ce58fe270f	2026-09-17 11:25:39.749708+00	2026-09-17 11:26:27.445569+00	2026-11-16 11:25:39.749494+00	2026-09-17 11:26:27.445443+00
9de7a899e6857cb803feb1e6e5a15244ec17ea339e44ab1ff41e90a2d9d2a9db	f829a7ff-9d6b-46f8-85e2-95680bde51e3	2026-09-17 11:26:39.32579+00	2026-09-17 11:26:39.32579+00	2026-11-16 11:26:39.325586+00	\N
\.


--
-- Data for Name: status_histories; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.status_histories (id, application_id, new_status_id, old_status_id, created_at) FROM stdin;
1	1	11	\N	2026-09-17 11:27:02.471238+00
2	2	11	\N	2026-09-17 11:27:27.227017+00
3	3	11	\N	2026-09-17 11:27:44.038616+00
4	4	11	\N	2026-09-17 11:27:54.724479+00
5	2	14	11	2026-09-17 11:28:20.708734+00
6	2	12	14	2026-09-17 11:28:48.117503+00
7	3	14	11	2026-09-17 11:29:17.597617+00
8	1	13	11	2026-09-17 11:29:44.270295+00
\.


--
-- Data for Name: user_settings; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.user_settings (id, created_at, updated_at, deleted_at, user_id, setup_step, default_application_status_id, is_banned, is_enabled, is_admin) FROM stdin;
1	2026-09-17 11:25:39.724902+00	2026-09-17 11:25:39.733583+00	\N	eba72c65-3e54-444a-8bad-d1ce58fe270f	1	2	f	t	t
2	2026-09-17 11:26:39.302381+00	2026-09-17 11:26:39.306687+00	\N	f829a7ff-9d6b-46f8-85e2-95680bde51e3	1	11	f	t	f
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.users (id, created_at, updated_at, deleted_at, email, username, password, user_settings_id) FROM stdin;
eba72c65-3e54-444a-8bad-d1ce58fe270f	2026-09-17 11:25:39.717117+00	2026-09-17 11:25:39.717117+00	\N	a@a.com	admin	$argon2id$v=19$m=65536,t=1,p=12$4o7F2ptei9M8jOdq3Q6Ktg$T6CwUqKJpe0VE4zIGyRBbA5nwv4VaYZK4vMlR5AYBKA	1
f829a7ff-9d6b-46f8-85e2-95680bde51e3	2026-09-17 11:26:39.294061+00	2026-09-17 11:26:39.294061+00	\N	test@test.com	test	$argon2id$v=19$m=65536,t=1,p=12$xDbVvyshonhZIdmHjf2Crg$lvvBHeu7WdKGsIoaDvxUmkUl2qbZabOgpxR/D/C6v7g	2
\.


--
-- Name: application_notes_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.application_notes_id_seq', 5, true);


--
-- Name: application_statuses_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.application_statuses_id_seq', 18, true);


--
-- Name: ban_lists_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.ban_lists_id_seq', 1, false);


--
-- Name: companies_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.companies_id_seq', 5, true);


--
-- Name: company_change_histories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.company_change_histories_id_seq', 1, false);


--
-- Name: job_applications_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.job_applications_id_seq', 4, true);


--
-- Name: status_histories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.status_histories_id_seq', 8, true);


--
-- Name: user_settings_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.user_settings_id_seq', 2, true);


--
-- Name: application_notes application_notes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.application_notes
    ADD CONSTRAINT application_notes_pkey PRIMARY KEY (id);


--
-- Name: application_statuses application_statuses_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.application_statuses
    ADD CONSTRAINT application_statuses_pkey PRIMARY KEY (id);


--
-- Name: ban_lists ban_lists_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ban_lists
    ADD CONSTRAINT ban_lists_pkey PRIMARY KEY (id);


--
-- Name: companies companies_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.companies
    ADD CONSTRAINT companies_pkey PRIMARY KEY (id);


--
-- Name: company_change_histories company_change_histories_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.company_change_histories
    ADD CONSTRAINT company_change_histories_pkey PRIMARY KEY (id);


--
-- Name: configs configs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.configs
    ADD CONSTRAINT configs_pkey PRIMARY KEY (key);


--
-- Name: job_applications job_applications_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.job_applications
    ADD CONSTRAINT job_applications_pkey PRIMARY KEY (id);


--
-- Name: refresh_tokens refresh_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_pkey PRIMARY KEY (token);


--
-- Name: status_histories status_histories_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.status_histories
    ADD CONSTRAINT status_histories_pkey PRIMARY KEY (id);


--
-- Name: user_settings user_settings_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_settings
    ADD CONSTRAINT user_settings_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_application_notes_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_application_notes_deleted_at ON public.application_notes USING btree (deleted_at);


--
-- Name: idx_application_statuses_archived; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_application_statuses_archived ON public.application_statuses USING btree (archived);


--
-- Name: idx_application_statuses_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_application_statuses_deleted_at ON public.application_statuses USING btree (deleted_at);


--
-- Name: idx_application_statuses_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_application_statuses_status ON public.application_statuses USING btree (status);


--
-- Name: idx_ban_lists_banned_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_ban_lists_banned_at ON public.ban_lists USING btree (banned_at);


--
-- Name: idx_ban_lists_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_ban_lists_deleted_at ON public.ban_lists USING btree (deleted_at);


--
-- Name: idx_companies_created_by; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_companies_created_by ON public.companies USING btree (created_by);


--
-- Name: idx_companies_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_companies_deleted_at ON public.companies USING btree (deleted_at);


--
-- Name: idx_companies_edited_by; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_companies_edited_by ON public.companies USING btree (edited_by);


--
-- Name: idx_companies_name; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_companies_name ON public.companies USING btree (name);


--
-- Name: idx_company_change_histories_changed_by; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_company_change_histories_changed_by ON public.company_change_histories USING btree (changed_by);


--
-- Name: idx_company_change_histories_company_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_company_change_histories_company_id ON public.company_change_histories USING btree (company_id);


--
-- Name: idx_configs_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_configs_deleted_at ON public.configs USING btree (deleted_at);


--
-- Name: idx_job_applications_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_job_applications_deleted_at ON public.job_applications USING btree (deleted_at);


--
-- Name: idx_refresh_tokens_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_refresh_tokens_user_id ON public.refresh_tokens USING btree (user_id);


--
-- Name: idx_user_settings_default_application_status_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_settings_default_application_status_id ON public.user_settings USING btree (default_application_status_id);


--
-- Name: idx_user_settings_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_settings_deleted_at ON public.user_settings USING btree (deleted_at);


--
-- Name: idx_users_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_deleted_at ON public.users USING btree (deleted_at);


--
-- Name: idx_users_email; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_users_email ON public.users USING btree (email);


--
-- Name: idx_users_user_settings_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_user_settings_id ON public.users USING btree (user_settings_id);


--
-- Name: idx_users_username; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_users_username ON public.users USING btree (username);


--
-- Name: application_notes fk_application_notes_status; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.application_notes
    ADD CONSTRAINT fk_application_notes_status FOREIGN KEY (status_id) REFERENCES public.application_statuses(id);


--
-- Name: ban_lists fk_ban_lists_banned_by_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ban_lists
    ADD CONSTRAINT fk_ban_lists_banned_by_user FOREIGN KEY (banned_by) REFERENCES public.users(id);


--
-- Name: companies fk_companies_created_by_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.companies
    ADD CONSTRAINT fk_companies_created_by_user FOREIGN KEY (created_by) REFERENCES public.users(id);


--
-- Name: companies fk_companies_edited_by_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.companies
    ADD CONSTRAINT fk_companies_edited_by_user FOREIGN KEY (edited_by) REFERENCES public.users(id);


--
-- Name: company_change_histories fk_company_change_histories_changed_by_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.company_change_histories
    ADD CONSTRAINT fk_company_change_histories_changed_by_user FOREIGN KEY (changed_by) REFERENCES public.users(id);


--
-- Name: company_change_histories fk_company_change_histories_company; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.company_change_histories
    ADD CONSTRAINT fk_company_change_histories_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- Name: application_notes fk_job_applications_notes; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.application_notes
    ADD CONSTRAINT fk_job_applications_notes FOREIGN KEY (application_id) REFERENCES public.job_applications(id);


--
-- Name: job_applications fk_job_applications_status; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.job_applications
    ADD CONSTRAINT fk_job_applications_status FOREIGN KEY (status_id) REFERENCES public.application_statuses(id);


--
-- Name: refresh_tokens fk_refresh_tokens_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT fk_refresh_tokens_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: status_histories fk_status_histories_application; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.status_histories
    ADD CONSTRAINT fk_status_histories_application FOREIGN KEY (application_id) REFERENCES public.job_applications(id);


--
-- Name: status_histories fk_status_histories_new_status; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.status_histories
    ADD CONSTRAINT fk_status_histories_new_status FOREIGN KEY (new_status_id) REFERENCES public.application_statuses(id);


--
-- Name: status_histories fk_status_histories_old_status; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.status_histories
    ADD CONSTRAINT fk_status_histories_old_status FOREIGN KEY (old_status_id) REFERENCES public.application_statuses(id);


--
-- Name: users fk_users_user_settings; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT fk_users_user_settings FOREIGN KEY (user_settings_id) REFERENCES public.user_settings(id);


--
-- PostgreSQL database dump complete
--

\unrestrict oEkl51gnaVYxkQZT9dBB3cDhbdVVjtb2OZSNKmn2pbLG93b1AtuKoVrXuM0avgE

