--
-- PostgreSQL database dump
--

-- Dumped from database version 12.2
-- Dumped by pg_dump version 12.2 (Ubuntu 12.2-2.pgdg18.04+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
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
-- Name: category; Type: TABLE; Schema: public; Owner: dsb
--

CREATE TABLE public.category (
    id integer NOT NULL,
    name character varying NOT NULL,
    "guildId" character varying NOT NULL
);


ALTER TABLE public.category OWNER TO dsb;

--
-- Name: category_id_seq; Type: SEQUENCE; Schema: public; Owner: dsb
--

CREATE SEQUENCE public.category_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.category_id_seq OWNER TO dsb;

--
-- Name: category_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: dsb
--

ALTER SEQUENCE public.category_id_seq OWNED BY public.category.id;


--
-- Name: image; Type: TABLE; Schema: public; Owner: dsb
--

CREATE TABLE public.image (
    id integer NOT NULL,
    uuid character varying NOT NULL
);


ALTER TABLE public.image OWNER TO dsb;

--
-- Name: image_id_seq; Type: SEQUENCE; Schema: public; Owner: dsb
--

CREATE SEQUENCE public.image_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.image_id_seq OWNER TO dsb;

--
-- Name: image_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: dsb
--

ALTER SEQUENCE public.image_id_seq OWNED BY public.image.id;


--
-- Name: sound; Type: TABLE; Schema: public; Owner: dsb
--

CREATE TABLE public.sound (
    id integer NOT NULL,
    uuid character varying NOT NULL,
    name character varying NOT NULL,
    "categoryId" integer,
    "guildId" character varying NOT NULL,
    "imageId" integer
);


ALTER TABLE public.sound OWNER TO dsb;

--
-- Name: sound_id_seq; Type: SEQUENCE; Schema: public; Owner: dsb
--

CREATE SEQUENCE public.sound_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.sound_id_seq OWNER TO dsb;

--
-- Name: sound_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: dsb
--

ALTER SEQUENCE public.sound_id_seq OWNED BY public.sound.id;


--
-- Name: category id; Type: DEFAULT; Schema: public; Owner: dsb
--

ALTER TABLE ONLY public.category ALTER COLUMN id SET DEFAULT nextval('public.category_id_seq'::regclass);


--
-- Name: image id; Type: DEFAULT; Schema: public; Owner: dsb
--

ALTER TABLE ONLY public.image ALTER COLUMN id SET DEFAULT nextval('public.image_id_seq'::regclass);


--
-- Name: sound id; Type: DEFAULT; Schema: public; Owner: dsb
--

ALTER TABLE ONLY public.sound ALTER COLUMN id SET DEFAULT nextval('public.sound_id_seq'::regclass);


--
-- Data for Name: category; Type: TABLE DATA; Schema: public; Owner: dsb
--

COPY public.category (id, name, "guildId") FROM stdin;
\.


--
-- Data for Name: image; Type: TABLE DATA; Schema: public; Owner: dsb
--

COPY public.image (id, uuid) FROM stdin;
1	6d110ae2-e38e-4a58-95ea-5f407d759133
2	be385c69-fe06-4f14-bddd-0a5cfbf7c5f8
3	32340f85-558e-4035-b77b-afe48ad27106
4	965e70e0-9152-432a-84b3-b09d4e09bc64
5	19ee7210-becb-47cf-bb8f-895a41114fa4
6	bdf9bc7b-2814-4154-b7a1-97a5eee7cb2c
7	5b3edc92-0f71-4a2a-b852-763245162bd1
8	ff7c07ef-1a01-40d4-8aad-86ddda9ebb04
9	3c03996e-0841-464d-9b23-70db2ec8d36f
\.


--
-- Data for Name: sound; Type: TABLE DATA; Schema: public; Owner: dsb
--

COPY public.sound (id, uuid, name, "categoryId", "guildId", "imageId") FROM stdin;
101	47559bc0-3734-436d-ad33-3bff28d1ac39	aled	\N	423590632582414346	8
106	ddb4ca3a-965a-4a21-98fe-59750b8e94c4	adedigato	\N	423590632582414346	9
\.


--
-- Name: category_id_seq; Type: SEQUENCE SET; Schema: public; Owner: dsb
--

SELECT pg_catalog.setval('public.category_id_seq', 49, true);


--
-- Name: image_id_seq; Type: SEQUENCE SET; Schema: public; Owner: dsb
--

SELECT pg_catalog.setval('public.image_id_seq', 9, true);


--
-- Name: sound_id_seq; Type: SEQUENCE SET; Schema: public; Owner: dsb
--

SELECT pg_catalog.setval('public.sound_id_seq', 106, true);


--
-- Name: sound PK_042a7f5e448107b2fd0eb4dfe8c; Type: CONSTRAINT; Schema: public; Owner: dsb
--

ALTER TABLE ONLY public.sound
    ADD CONSTRAINT "PK_042a7f5e448107b2fd0eb4dfe8c" PRIMARY KEY (id);


--
-- Name: category PK_9c4e4a89e3674fc9f382d733f03; Type: CONSTRAINT; Schema: public; Owner: dsb
--

ALTER TABLE ONLY public.category
    ADD CONSTRAINT "PK_9c4e4a89e3674fc9f382d733f03" PRIMARY KEY (id);


--
-- Name: image PK_d6db1ab4ee9ad9dbe86c64e4cc3; Type: CONSTRAINT; Schema: public; Owner: dsb
--

ALTER TABLE ONLY public.image
    ADD CONSTRAINT "PK_d6db1ab4ee9ad9dbe86c64e4cc3" PRIMARY KEY (id);


--
-- Name: sound UQ_0a01fdbad724d02f97b02712497; Type: CONSTRAINT; Schema: public; Owner: dsb
--

ALTER TABLE ONLY public.sound
    ADD CONSTRAINT "UQ_0a01fdbad724d02f97b02712497" UNIQUE ("imageId");


--
-- Name: sound FK_0a01fdbad724d02f97b02712497; Type: FK CONSTRAINT; Schema: public; Owner: dsb
--

ALTER TABLE ONLY public.sound
    ADD CONSTRAINT "FK_0a01fdbad724d02f97b02712497" FOREIGN KEY ("imageId") REFERENCES public.image(id);


--
-- Name: sound FK_496fd8fdeaace05755caddad1da; Type: FK CONSTRAINT; Schema: public; Owner: dsb
--

ALTER TABLE ONLY public.sound
    ADD CONSTRAINT "FK_496fd8fdeaace05755caddad1da" FOREIGN KEY ("categoryId") REFERENCES public.category(id) ON DELETE SET NULL;


--
-- PostgreSQL database dump complete
--

