"use strict";

const $ = (id) => document.getElementById(id);
const errorBox = $("errorBox");

function showError(msg) {
  errorBox.hidden = false;
  errorBox.textContent = msg;
}

function clearError() {
  errorBox.hidden = true;
  errorBox.textContent = "";
}

async function postJSON(url, body) {
  clearError();
  const res = await fetch(url, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  let data = null;
  try {
    data = await res.json();
  } catch (e) {
    data = {};
  }
  if (!res.ok) {
    throw new Error((data && data.error) ? data.error : ("HTTP " + res.status));
  }
  return data;
}

function fmt(v) {
  if (!isFinite(v)) {
    return String(v);
  }
  return v.toExponential(4);
}

function fmtPct(v) {
  return (v * 100).toFixed(2) + "%";
}

function readInput() {
  return {
    r: Number($("r").value),
    gamma: Number($("gamma").value),
    intake: { t: Number($("t1").value), p: Number($("p1").value) },
    qin: Number($("qin").value),
    mode: $("mode").value,
  };
}

function fillInput(data) {
  $("r").value = data.r;
  $("gamma").value = data.gamma;
  $("t1").value = data.intake.t;
  $("p1").value = data.intake.p;
  $("qin").value = data.qin;
  $("mode").value = data.mode || "otto";
}

$("loadExample").addEventListener("click", async () => {
  const res = await fetch("/api/examples");
  if (!res.ok) {
    showError("failed to load examples from backend");
    return;
  }
  const all = await res.json();
  const data = all.r8;
  if (!data) {
    showError("example r8 not found");
    return;
  }
  fillInput(data);
  compute();
});

$("computeBtn").addEventListener("click", compute);

async function compute() {
  const body = readInput();
  try {
    const cyc = await postJSON("/api/cycle", body);
    renderCycle(cyc);
    const pv = await postJSON("/api/pv", body);
    drawPV(pv);
  } catch (e) {
    showError(e.message);
  }
}

function renderCycle(c) {
  const rows = c.states.map((s) =>
    `<tr><td>${s.point}</td><td>${fmt(s.p)}</td><td>${fmt(s.v)}</td><td>${fmt(s.t)}</td></tr>`
  ).join("");
  $("stateBody").innerHTML = rows;
  $("metrics").innerHTML =
    `<p><strong>η = ${fmtPct(c.eta)}</strong>（闭式 1−r^{1−γ} = ${fmtPct(c.eta_close)}，w_net/q_in = ${fmtPct(c.eta_work)}）</p>` +
    `<p>w_net = ${fmt(c.w_net)} J/kg · q_out = ${fmt(c.q_out)} J/kg · q_in = ${fmt(c.q_in)} J/kg</p>` +
    `<p>MEP = ${fmt(c.mep)} Pa · 循环排量 v1−v2 = ${fmt(c.states[0].v - c.states[1].v)} m³/kg</p>` +
    `<p>p–v 闭合面积（数值积分） = ${fmt(c.area)} J/kg</p>`;
}

function drawPV(pv) {
  const svg = $("pvSvg");
  svg.innerHTML = "";
  const curve = pv.curve;
  if (!curve || curve.length < 2) {
    $("pvNote").textContent = "curve empty";
    return;
  }
  let minV = Infinity, maxV = -Infinity, minP = Infinity, maxP = -Infinity;
  for (const c of curve) {
    minV = Math.min(minV, c.v);
    maxV = Math.max(maxV, c.v);
    minP = Math.min(minP, c.p);
    maxP = Math.max(maxP, c.p);
  }
  const pad = 14;
  const W = 560, H = 340;
  const spanV = maxV - minV;
  const spanLogP = Math.log(maxP) - Math.log(minP);
  const x = (v) => pad + ((v - minV) / spanV) * (W - 2 * pad);
  const y = (p) => H - pad - ((Math.log(p) - Math.log(minP)) / spanLogP) * (H - 2 * pad);

  const points = curve.map((c) => `${x(c.v)},${y(c.p)}`).join(" ");
  const path = document.createElementNS("http://www.w3.org/2000/svg", "polyline");
  path.setAttribute("points", points);
  path.setAttribute("fill", "none");
  path.setAttribute("stroke", "#1a5fb4");
  path.setAttribute("stroke-width", "2");
  svg.appendChild(path);

  const first = curve[0];
  const marker = document.createElementNS("http://www.w3.org/2000/svg", "circle");
  marker.setAttribute("cx", x(first.v));
  marker.setAttribute("cy", y(first.p));
  marker.setAttribute("r", "4");
  marker.setAttribute("fill", "#e01b24");
  svg.appendChild(marker);
  const label = document.createElementNS("http://www.w3.org/2000/svg", "text");
  label.setAttribute("x", x(first.v) + 6);
  label.setAttribute("y", y(first.p) - 6);
  label.setAttribute("fill", "#e01b24");
  label.textContent = "1";
  svg.appendChild(label);

  const axisX = document.createElementNS("http://www.w3.org/2000/svg", "line");
  axisX.setAttribute("x1", pad);
  axisX.setAttribute("y1", H - pad);
  axisX.setAttribute("x2", W - pad);
  axisX.setAttribute("y2", H - pad);
  axisX.setAttribute("stroke", "#888");
  svg.appendChild(axisX);
  const axisY = document.createElementNS("http://www.w3.org/2000/svg", "line");
  axisY.setAttribute("x1", pad);
  axisY.setAttribute("y1", pad);
  axisY.setAttribute("x2", pad);
  axisY.setAttribute("y2", H - pad);
  axisY.setAttribute("stroke", "#888");
  svg.appendChild(axisY);

  $("pvNote").textContent =
    `mode=${pv.mode} · area=${fmt(pv.area)} J/kg · w_net=${fmt(pv.w_net)} J/kg · η=${fmtPct(pv.eta)}（y 轴对数刻度）`;
}
