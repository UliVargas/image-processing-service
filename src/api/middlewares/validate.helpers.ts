import type { RequestHandler } from "express";
import type * as v from "valibot";
import { validate } from "./validate.middleware";

export const validateBody = <T extends v.GenericSchema>(
	schema: T,
): RequestHandler => validate({ schema, type: "body" });

export const validateParams = <T extends v.GenericSchema>(
	schema: T,
): RequestHandler => validate({ schema, type: "params" });

export const validateQuery = <T extends v.GenericSchema>(
	schema: T,
): RequestHandler => validate({ schema, type: "query" });
