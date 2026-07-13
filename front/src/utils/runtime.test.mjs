import test from "node:test";
import assert from "node:assert/strict";
import { normalizeApiBase } from "./runtime.js";

test("接口地址为空时使用同源路径", () => {
  assert.equal(normalizeApiBase(""), "");
  assert.equal(normalizeApiBase(null), "");
});

test("接口地址会去掉末尾斜杠", () => {
  assert.equal(normalizeApiBase("https://example.test///"), "https://example.test");
});
