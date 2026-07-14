import test from "node:test";
import assert from "node:assert/strict";

import { apiBase, normalizeApiBase, parseAPIResponse } from "./runtime.js";

test("默认接口地址指向本机网关", () => {
  assert.equal(apiBase, "http://127.0.0.1:8080");
});

test("接口地址会去掉末尾斜杠", () => {
  assert.equal(normalizeApiBase("https://example.test///"), "https://example.test");
});

test("非结构化响应返回明确的网关地址错误", async () => {
  const response = new Response("<html>开发服务器页面</html>", {
    headers: { "content-type": "text/html; charset=utf-8" },
  });

  await assert.rejects(
    () => parseAPIResponse(response),
    /网关返回了非接口数据，请检查网关地址/
  );
});

test("结构化响应返回解析后的数据", async () => {
  const response = new Response(JSON.stringify({ code: 0, msg: "ok" }), {
    headers: { "content-type": "application/json; charset=utf-8" },
  });

  assert.deepEqual(await parseAPIResponse(response), { code: 0, msg: "ok" });
});
