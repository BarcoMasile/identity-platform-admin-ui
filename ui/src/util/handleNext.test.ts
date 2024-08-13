import { handleNext } from "./handleNext";

test("handles the 'next' path param", () => {
  window.location.href = "/old/?next=/new";
  handleNext();
  expect(window.location.pathname).toBe("/new");
});

test("handles the 'next' param with domain", () => {
  window.location.href = "http://example.com/old/?next=/new";
  handleNext();
  expect(window.location.pathname).toBe("/new");
});

test("handles no 'next' param", () => {
  window.location.href = "/current/?search=query";
  handleNext();
  expect(window.location.pathname).toBe("/current/");
  expect(window.location.search).toBe("?search=query");
});

test("no redirect if the next param matches the current page", () => {
  window.location.href = "/current/?next=/current";
  handleNext();
  expect(window.location.pathname).toBe("/current/");
  expect(window.location.search).toBe("");
});

test("no redirect if the next param has a different domain", () => {
  window.location.href =
    "http://example.com/old/?next=http://notexample.com/new";
  handleNext();
  expect(window.location.host).toBe("example.com");
  expect(window.location.pathname).toBe("/old/");
  expect(window.location.search).toBe("");
});
