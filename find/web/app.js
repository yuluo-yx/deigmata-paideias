const statusEl = document.getElementById("status");
const startBtn = document.getElementById("startBtn");
const domainInput = document.getElementById("domainInput");
const depthInput = document.getElementById("depthInput");
const workersInput = document.getElementById("workersInput");
const treeContainer = document.getElementById("treeContainer");

let currentScanId = null;
let pollTimer = null;

function setStatus(text, mode = "idle") {
  statusEl.textContent = text;
  statusEl.className = `status status-${mode}`;
}

function renderTree(node) {
  if (!node) {
    treeContainer.innerHTML = "<div class=\"empty\">暂无数据</div>";
    return;
  }
  treeContainer.innerHTML = "";
  const layout = buildLayout(node);
  const svg = renderSvg(layout);
  treeContainer.appendChild(svg);
}

function buildLayout(root) {
  const xGap = 240;
  const yGap = 110;
  let leafIndex = 0;

  const nodes = [];
  const links = [];

  function traverse(node, depth) {
    const current = {
      id: nodes.length,
      node,
      depth,
      x: depth * xGap,
      y: 0,
    };
    nodes.push(current);

    if (!node.children || node.children.length === 0) {
      current.y = leafIndex * yGap;
      leafIndex += 1;
    } else {
      const childLayouts = node.children.map((child) => traverse(child, depth + 1));
      const minY = Math.min(...childLayouts.map((c) => c.y));
      const maxY = Math.max(...childLayouts.map((c) => c.y));
      current.y = (minY + maxY) / 2;
      childLayouts.forEach((childLayout) => {
        links.push({ from: current, to: childLayout });
      });
    }

    return current;
  }

  traverse(root, 0);

  const width = Math.max(600, (Math.max(...nodes.map((n) => n.x)) || 0) + 320);
  const height = Math.max(300, (Math.max(...nodes.map((n) => n.y)) || 0) + 140);

  return { nodes, links, width, height };
}

function renderSvg(layout) {
  const svg = document.createElementNS("http://www.w3.org/2000/svg", "svg");
  svg.setAttribute("class", "tree-svg");
  svg.setAttribute("width", layout.width);
  svg.setAttribute("height", layout.height);
  svg.setAttribute("viewBox", `0 0 ${layout.width} ${layout.height}`);

  const linkGroup = document.createElementNS("http://www.w3.org/2000/svg", "g");
  linkGroup.setAttribute("class", "link-group");

  layout.links.forEach(({ from, to }) => {
    const path = document.createElementNS("http://www.w3.org/2000/svg", "path");
    const startX = from.x + 180;
    const startY = from.y + 40;
    const endX = to.x + 10;
    const endY = to.y + 40;
    const midX = (startX + endX) / 2;
    path.setAttribute(
      "d",
      `M ${startX} ${startY} C ${midX} ${startY}, ${midX} ${endY}, ${endX} ${endY}`
    );
    path.setAttribute("class", "tree-link");
    linkGroup.appendChild(path);
  });

  const nodeGroup = document.createElementNS("http://www.w3.org/2000/svg", "g");
  nodeGroup.setAttribute("class", "node-group");

  layout.nodes.forEach(({ node, x, y }) => {
    const group = document.createElementNS("http://www.w3.org/2000/svg", "g");
    group.setAttribute("class", "tree-node");
    group.setAttribute("transform", `translate(${x}, ${y})`);

    const rect = document.createElementNS("http://www.w3.org/2000/svg", "rect");
    rect.setAttribute("width", "190");
    rect.setAttribute("height", "80");
    rect.setAttribute("rx", "12");
    rect.setAttribute("ry", "12");
    rect.setAttribute("class", "node-rect");

    const title = document.createElementNS("http://www.w3.org/2000/svg", "text");
    title.setAttribute("x", "12");
    title.setAttribute("y", "24");
    title.setAttribute("class", "node-title");
    title.textContent = node.domain;

    const meta = document.createElementNS("http://www.w3.org/2000/svg", "text");
    meta.setAttribute("x", "12");
    meta.setAttribute("y", "44");
    meta.setAttribute("class", "node-meta");

    const records = node.records && node.records.length > 0 ? node.records : [];
    if (records.length === 0) {
      const tspan = document.createElementNS("http://www.w3.org/2000/svg", "tspan");
      tspan.setAttribute("x", "12");
      tspan.setAttribute("dy", "0");
      tspan.textContent = "无解析记录";
      meta.appendChild(tspan);
    } else {
      records.slice(0, 2).forEach((record, index) => {
        const tspan = document.createElementNS("http://www.w3.org/2000/svg", "tspan");
        tspan.setAttribute("x", "12");
        tspan.setAttribute("dy", index === 0 ? "0" : "16");
        tspan.textContent = `${record.ip} ${record.location || "未知"}`;
        meta.appendChild(tspan);
      });
      if (records.length > 2) {
        const tspan = document.createElementNS("http://www.w3.org/2000/svg", "tspan");
        tspan.setAttribute("x", "12");
        tspan.setAttribute("dy", "16");
        tspan.textContent = `+${records.length - 2} 更多`;
        meta.appendChild(tspan);
      }
    }

    group.appendChild(rect);
    group.appendChild(title);
    group.appendChild(meta);
    nodeGroup.appendChild(group);
  });

  svg.appendChild(linkGroup);
  svg.appendChild(nodeGroup);
  return svg;
}

async function startScan() {
  const domain = domainInput.value.trim();
  const maxDepth = Number(depthInput.value);
  const workers = Number(workersInput.value);

  if (!domain) {
    setStatus("请输入根域名", "error");
    return;
  }

  setStatus("正在提交扫描请求...", "running");

  try {
    const response = await fetch("/scan", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ domain, maxDepth, workers }),
    });

    if (!response.ok) {
      throw new Error("提交扫描失败");
    }

    const data = await response.json();
    currentScanId = data.scanId;
    setStatus("扫描中...", "running");

    if (pollTimer) {
      clearInterval(pollTimer);
    }

    pollTimer = setInterval(fetchTree, 2000);
    await fetchTree();
  } catch (error) {
    setStatus(error.message || "扫描启动失败", "error");
  }
}

async function fetchTree() {
  if (!currentScanId) return;

  try {
    const response = await fetch(`/tree/${currentScanId}`);
    if (!response.ok) {
      throw new Error("获取树数据失败");
    }

    const data = await response.json();
    renderTree(data.tree);

    if (data.status === "finished") {
      setStatus("扫描完成", "success");
      clearInterval(pollTimer);
      pollTimer = null;
    } else if (data.status === "error") {
      setStatus(data.error || "扫描出错", "error");
      clearInterval(pollTimer);
      pollTimer = null;
    } else {
      setStatus("扫描中...", "running");
    }
  } catch (error) {
    setStatus(error.message || "获取树失败", "error");
  }
}

startBtn.addEventListener("click", startScan);

renderTree(null);
