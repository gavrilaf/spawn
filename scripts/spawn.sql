--
-- PostgreSQL database dump
--

-- Dumped from database version 10.0
-- Dumped by pg_dump version 10.0

-- Started on 2017-12-02 22:44:08 EET

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

DROP DATABASE spawn;
--
-- TOC entry 2893 (class 1262 OID 24576)
-- Name: spawn; Type: DATABASE; Schema: -; Owner: postgres
--

CREATE DATABASE spawn WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'en_US.utf8' LC_CTYPE = 'en_US.utf8';


ALTER DATABASE spawn OWNER TO postgres;

\connect spawn

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 1 (class 3079 OID 12980)
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;

--
-- TOC entry 2896 (class 0 OID 0)
-- Dependencies: 1
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- TOC entry 199 (class 1259 OID 24622)
-- Name: Clients; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE "Clients" (
    id text NOT NULL,
    secret text NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    description text DEFAULT ''::text NOT NULL,
    def_scope bigint DEFAULT 0 NOT NULL
);


ALTER TABLE "Clients" OWNER TO postgres;

--
-- TOC entry 196 (class 1259 OID 24577)
-- Name: Devices; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE "Devices" (
    device_id text NOT NULL,
    device_name text NOT NULL,
    user_id uuid NOT NULL,
    is_confirmed boolean DEFAULT false NOT NULL,
    fingerprint text DEFAULT ''::text NOT NULL,
    locale text DEFAULT 'en'::text NOT NULL,
    lang text DEFAULT 'en'::text NOT NULL
);


ALTER TABLE "Devices" OWNER TO postgres;

--
-- TOC entry 197 (class 1259 OID 24587)
-- Name: LoginsLog; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE "LoginsLog" (
    user_id uuid NOT NULL,
    device_id text NOT NULL,
    device_name text NOT NULL,
    "timestamp" time with time zone NOT NULL,
    user_agent text DEFAULT ''::text NOT NULL,
    ip text DEFAULT ''::text NOT NULL,
    region text DEFAULT ''::text NOT NULL
);


ALTER TABLE "LoginsLog" OWNER TO postgres;

--
-- TOC entry 198 (class 1259 OID 24596)
-- Name: Users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE "Users" (
    id uuid NOT NULL,
    username text NOT NULL,
    password text NOT NULL,
    is_locked boolean DEFAULT false NOT NULL,
    is_email_confirmed boolean DEFAULT false NOT NULL,
    is_2fa_required boolean DEFAULT false NOT NULL,
    scope bigint DEFAULT 0 NOT NULL,
    first_name text DEFAULT ''::text NOT NULL,
    last_name text DEFAULT ''::text NOT NULL,
    birth_date date DEFAULT '1800-01-01'::date NOT NULL,
    country text DEFAULT ''::text NOT NULL,
    phone_country_code integer DEFAULT 0 NOT NULL,
    phone_number text DEFAULT ''::text NOT NULL,
    is_phone_confirmed boolean DEFAULT false NOT NULL
);


ALTER TABLE "Users" OWNER TO postgres;

--
-- TOC entry 2750 (class 2606 OID 16415)
-- Name: Users Users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY "Users"
    ADD CONSTRAINT "Users_pkey" PRIMARY KEY (id);


--
-- TOC entry 2752 (class 2606 OID 16417)
-- Name: Users username_uniq; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY "Users"
    ADD CONSTRAINT username_uniq UNIQUE (username);


--
-- TOC entry 2753 (class 2606 OID 16418)
-- Name: Devices device_user_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY "Devices"
    ADD CONSTRAINT device_user_fk FOREIGN KEY (user_id) REFERENCES "Users"(id);

--
-- TOC entry 2888 (class 0 OID 24622)
-- Dependencies: 199
-- Data for Name: Clients; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO "Clients" VALUES ('ios-client-01-key', 'sA?2,S]$P6''Cs`Q)&4;18LXIj#b_=D', true, 'Spawn iOS application key', 0);
INSERT INTO "Clients" VALUES ('client-test-01', '~_7|cjU^L?l5JI/jqN)S7|-I;=wz6<', true, 'Client for internal & external testing', 0);


--
-- TOC entry 2885 (class 0 OID 24577)
-- Dependencies: 196
-- Data for Name: Devices; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- TOC entry 2886 (class 0 OID 24587)
-- Dependencies: 197
-- Data for Name: LoginsLog; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- TOC entry 2887 (class 0 OID 24596)
-- Dependencies: 198
-- Data for Name: Users; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- TOC entry 2895 (class 0 OID 0)
-- Dependencies: 5
-- Name: public; Type: ACL; Schema: -; Owner: postgres
--

GRANT ALL ON SCHEMA public TO PUBLIC;


-- Completed on 2017-12-02 22:44:08 EET

--
-- PostgreSQL database dump complete
--

