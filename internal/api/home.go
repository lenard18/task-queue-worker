package api

// homeHTML es la página de bienvenida visual del Task Queue Worker
// Se sirve en GET / para que se vea bien en el navegador
const homeHTML = `<!DOCTYPE html>
<html lang="es">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Task Queue Worker — Leonardo Gómez</title>
<style>
  *{box-sizing:border-box;margin:0;padding:0}
  :root{
    --bg:#0f172a;--bg2:#1e293b;--bg3:#334155;
    --go:#00ADD8;--go2:#007d9c;
    --green:#22c55e;--amber:#f59e0b;--red:#ef4444;--blue:#3b82f6;
    --text:#f1f5f9;--muted:#94a3b8;--border:#334155;
  }
  body{font-family:'Segoe UI',Arial,sans-serif;background:var(--bg);color:var(--text);min-height:100vh}

  /* ── HEADER ── */
  header{
    background:linear-gradient(135deg,#003B57 0%,#00ADD8 100%);
    padding:2.5rem 1.5rem 2rem;text-align:center;position:relative;overflow:hidden;
  }
  header::before{
    content:'';position:absolute;inset:0;
    background:url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23ffffff' fill-opacity='0.05'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E");
  }
  .header-badge{
    display:inline-block;background:rgba(255,255,255,0.15);
    color:#fff;font-size:0.7rem;font-weight:700;letter-spacing:0.1em;
    text-transform:uppercase;padding:0.3rem 0.9rem;border-radius:20px;
    border:1px solid rgba(255,255,255,0.25);margin-bottom:1rem;
  }
  header h1{font-size:2rem;font-weight:800;color:#fff;margin-bottom:0.4rem}
  header h1 span{color:#7dd3fc}
  header p{color:rgba(255,255,255,0.8);font-size:0.95rem}
  .author{margin-top:1rem;font-size:0.8rem;color:rgba(255,255,255,0.65)}
  .author a{color:#7dd3fc;text-decoration:none}

  /* ── LAYOUT ── */
  .container{max-width:900px;margin:0 auto;padding:2rem 1.5rem}

  /* ── STATS ── */
  .stats-grid{display:grid;grid-template-columns:repeat(4,1fr);gap:0.75rem;margin-bottom:2rem}
  @media(max-width:600px){.stats-grid{grid-template-columns:repeat(2,1fr)}}
  .stat-card{
    background:var(--bg2);border:1px solid var(--border);border-radius:12px;
    padding:1.1rem;text-align:center;transition:all 0.2s;
  }
  .stat-card:hover{border-color:var(--go);transform:translateY(-2px)}
  .stat-num{font-size:2rem;font-weight:800;margin-bottom:0.2rem}
  .stat-num.green{color:var(--green)}
  .stat-num.amber{color:var(--amber)}
  .stat-num.red{color:var(--red)}
  .stat-num.blue{color:var(--blue)}
  .stat-label{font-size:0.7rem;font-weight:600;letter-spacing:0.08em;text-transform:uppercase;color:var(--muted)}

  /* ── PANEL DEMO ── */
  .demo-panel{
    background:var(--bg2);border:1px solid var(--border);border-radius:16px;
    padding:1.5rem;margin-bottom:2rem;
  }
  .panel-title{
    font-size:0.72rem;font-weight:700;letter-spacing:0.1em;text-transform:uppercase;
    color:var(--muted);margin-bottom:1rem;display:flex;align-items:center;gap:0.5rem;
  }
  .btn{
    display:inline-flex;align-items:center;gap:0.4rem;
    padding:0.6rem 1.2rem;border-radius:8px;font-size:0.85rem;font-weight:700;
    border:none;cursor:pointer;transition:all 0.15s;text-decoration:none;
  }
  .btn-go{background:var(--go);color:#fff}
  .btn-go:hover{background:var(--go2);transform:translateY(-1px)}
  .btn-outline{background:transparent;color:var(--text);border:1.5px solid var(--border)}
  .btn-outline:hover{border-color:var(--go);color:var(--go)}
  .btn-group{display:flex;flex-wrap:wrap;gap:0.6rem;margin-bottom:1.2rem}

  /* ── LOG DE JOBS ── */
  .jobs-log{
    background:var(--bg);border:1px solid var(--border);border-radius:10px;
    padding:1rem;max-height:260px;overflow-y:auto;font-family:monospace;font-size:0.8rem;
  }
  .log-empty{color:var(--muted);text-align:center;padding:1rem}
  .job-row{
    display:flex;align-items:center;gap:0.6rem;padding:0.4rem 0;
    border-bottom:1px solid rgba(255,255,255,0.04);
  }
  .job-row:last-child{border-bottom:none}
  .job-id{color:var(--muted);min-width:28px}
  .job-type{color:#7dd3fc;flex:1}
  .job-status{
    font-size:0.68rem;font-weight:700;padding:0.15rem 0.5rem;border-radius:5px;
    text-transform:uppercase;white-space:nowrap;
  }
  .status-done{background:#14532d;color:var(--green)}
  .status-pending{background:#1c1917;color:var(--amber)}
  .status-processing{background:#1e3a5f;color:var(--blue)}
  .status-failed{background:#450a0a;color:var(--red)}
  .job-pri{color:var(--muted);font-size:0.7rem;min-width:40px;text-align:right}

  /* ── ENDPOINTS ── */
  .endpoints{display:grid;gap:0.5rem;margin-bottom:2rem}
  .endpoint-row{
    background:var(--bg2);border:1px solid var(--border);border-radius:10px;
    padding:0.75rem 1rem;display:flex;align-items:center;gap:0.75rem;
  }
  .method{
    font-size:0.68rem;font-weight:800;padding:0.2rem 0.5rem;border-radius:5px;
    min-width:42px;text-align:center;font-family:monospace;
  }
  .method.GET{background:#1e3a5f;color:#60a5fa}
  .method.POST{background:#14532d;color:#4ade80}
  .endpoint-path{font-family:monospace;font-size:0.85rem;color:var(--go);flex:1}
  .endpoint-desc{font-size:0.78rem;color:var(--muted)}

  /* ── TECH STACK ── */
  .tech-grid{display:flex;flex-wrap:wrap;gap:0.5rem;margin-bottom:2rem}
  .tech-pill{
    font-size:0.75rem;font-weight:600;padding:0.3rem 0.75rem;border-radius:8px;
    background:var(--bg2);border:1px solid var(--border);color:var(--text);
    font-family:monospace;
  }
  .tech-pill.go{background:#003B57;border-color:var(--go);color:#7dd3fc}
  .tech-pill.pg{background:#1a2e4a;border-color:#336791;color:#79b4e8}

  /* ── FOOTER ── */
  footer{
    text-align:center;padding:2rem;color:var(--muted);font-size:0.8rem;
    border-top:1px solid var(--border);
  }
  footer a{color:var(--go);text-decoration:none}

  /* ── SPINNER ── */
  .spin{animation:spin 1s linear infinite;display:inline-block}
  @keyframes spin{to{transform:rotate(360deg)}}

  /* ── ALERT ── */
  .alert{
    padding:0.75rem 1rem;border-radius:8px;font-size:0.82rem;margin-top:0.75rem;
    display:none;
  }
  .alert.show{display:block}
  .alert-ok{background:#14532d;color:#4ade80;border:1px solid #166534}
  .alert-err{background:#450a0a;color:#f87171;border:1px solid #7f1d1d}
</style>
</head>
<body>

<!-- HEADER -->
<header>
  <div class="header-badge">&#x1F504; Task Queue Worker &mdash; En Producci&oacute;n</div>
  <h1>Task Queue <span>Worker</span></h1>
  <p>Sistema de procesamiento as&iacute;ncrono con workers concurrentes</p>
  <div class="author">
    Desarrollado por
    <a href="https://github.com/lenard18" target="_blank">Leonardo G&oacute;mez</a>
    &middot; Go + PostgreSQL + Docker
  </div>
</header>

<div class="container">

  <!-- STATS EN VIVO -->
  <div class="panel-title">&#x1F4CA; Estad&iacute;sticas de la cola &mdash; <span id="last-update">cargando...</span></div>
  <div class="stats-grid">
    <div class="stat-card">
      <div class="stat-num green" id="stat-done">—</div>
      <div class="stat-label">&#x2705; Completados</div>
    </div>
    <div class="stat-card">
      <div class="stat-num amber" id="stat-pending">—</div>
      <div class="stat-label">&#x23F3; Pendientes</div>
    </div>
    <div class="stat-card">
      <div class="stat-num blue" id="stat-processing">—</div>
      <div class="stat-label">&#x2699;&#xFE0F; Procesando</div>
    </div>
    <div class="stat-card">
      <div class="stat-num red" id="stat-failed">—</div>
      <div class="stat-label">&#x274C; Fallidos</div>
    </div>
  </div>

  <!-- PANEL DEMO -->
  <div class="demo-panel">
    <div class="panel-title">&#x1F528; Probar el sistema en vivo</div>
    <div class="btn-group">
      <button class="btn btn-go" onclick="runDemo()">
        &#x25B6; Encolar 5 jobs de demo
      </button>
      <button class="btn btn-outline" onclick="loadJobs()">
        &#x1F504; Actualizar jobs
      </button>
      <button class="btn btn-outline" onclick="enqueueAlert()">
        &#x1F6A8; Enviar alerta cr&iacute;tica
      </button>
    </div>
    <div id="demo-alert" class="alert"></div>
    <div class="jobs-log" id="jobs-log">
      <div class="log-empty">Haz clic en "Encolar 5 jobs de demo" para ver el sistema en acci&oacute;n</div>
    </div>
  </div>

  <!-- ENDPOINTS -->
  <div class="panel-title">&#x1F517; Endpoints disponibles</div>
  <div class="endpoints">
    <div class="endpoint-row">
      <span class="method POST">POST</span>
      <span class="endpoint-path">/api/jobs</span>
      <span class="endpoint-desc">Encolar un nuevo job (send_email, send_welcome, send_report, send_alert)</span>
    </div>
    <div class="endpoint-row">
      <span class="method GET">GET</span>
      <span class="endpoint-path">/api/jobs</span>
      <span class="endpoint-desc">Listar jobs &mdash; filtros: ?status=pending|done|failed&amp;limit=50</span>
    </div>
    <div class="endpoint-row">
      <span class="method GET">GET</span>
      <span class="endpoint-path">/api/stats</span>
      <span class="endpoint-desc">Dashboard de estad&iacute;sticas de la cola</span>
    </div>
    <div class="endpoint-row">
      <span class="method POST">POST</span>
      <span class="endpoint-path">/api/demo</span>
      <span class="endpoint-desc">Encolar 5 jobs de ejemplo para ver el sistema en acci&oacute;n</span>
    </div>
    <div class="endpoint-row">
      <span class="method GET">GET</span>
      <span class="endpoint-path">/health</span>
      <span class="endpoint-desc">Health check &mdash; verifica que el servicio est&eacute; activo</span>
    </div>
  </div>

  <!-- TECH STACK -->
  <div class="panel-title">&#x26A1; Stack tecnol&oacute;gico</div>
  <div class="tech-grid">
    <span class="tech-pill go">Go 1.22</span>
    <span class="tech-pill pg">PostgreSQL 16</span>
    <span class="tech-pill">Gin Framework</span>
    <span class="tech-pill">Goroutines</span>
    <span class="tech-pill">Channels</span>
    <span class="tech-pill">SELECT FOR UPDATE SKIP LOCKED</span>
    <span class="tech-pill">Backoff exponencial</span>
    <span class="tech-pill">SMTP net/smtp</span>
    <span class="tech-pill">Docker multi-stage</span>
    <span class="tech-pill">Render.com</span>
  </div>

</div>

<footer>
  <p>
    <a href="https://github.com/lenard18/task-queue-worker" target="_blank">&#x2325; Ver c&oacute;digo en GitHub</a>
    &nbsp;&middot;&nbsp;
    <a href="https://lenard18.github.io" target="_blank">&#x1F464; Portfolio completo</a>
    &nbsp;&middot;&nbsp;
    Desarrollado por Leonardo G&oacute;mez
  </p>
</footer>

<script>
const BASE = '';

async function loadStats() {
  try {
    const r = await fetch(BASE + '/api/stats');
    const d = await r.json();
    const q = d.queue_stats || {};
    document.getElementById('stat-done').textContent       = (q.done       || {}).total ?? 0;
    document.getElementById('stat-pending').textContent    = (q.pending    || {}).total ?? 0;
    document.getElementById('stat-processing').textContent = (q.processing || {}).total ?? 0;
    document.getElementById('stat-failed').textContent     = (q.failed     || {}).total ?? 0;
    document.getElementById('last-update').textContent     = 'actualizado ' + new Date().toLocaleTimeString('es-CO');
  } catch(e) {
    document.getElementById('last-update').textContent = 'error al cargar';
  }
}

async function loadJobs() {
  try {
    const r = await fetch(BASE + '/api/jobs?limit=20');
    const d = await r.json();
    renderJobs(d.jobs || []);
    loadStats();
  } catch(e) {
    showAlert('Error al cargar jobs: ' + e.message, false);
  }
}

function renderJobs(jobs) {
  const el = document.getElementById('jobs-log');
  if (!jobs.length) {
    el.innerHTML = '<div class="log-empty">No hay jobs todav&iacute;a &mdash; haz clic en "Encolar 5 jobs de demo"</div>';
    return;
  }
  el.innerHTML = jobs.map(j => {
    const statusClass = 'status-' + j.status;
    const statusEmoji = {done:'✅',pending:'⏳',processing:'⚙️',failed:'❌'}[j.status] || '•';
    return '<div class="job-row">'
      + '<span class="job-id">#' + j.id + '</span>'
      + '<span class="job-type">' + j.type + '</span>'
      + '<span class="job-status ' + statusClass + '">' + statusEmoji + ' ' + j.status + '</span>'
      + '<span class="job-pri">p=' + j.priority + '</span>'
      + '</div>';
  }).join('');
}

async function runDemo() {
  const btn = event.target;
  btn.disabled = true;
  btn.innerHTML = '<span class="spin">&#x1F504;</span> Encolando...';
  try {
    const r = await fetch(BASE + '/api/demo', {method:'POST'});
    const d = await r.json();
    showAlert('&#x2705; ' + d.message, true);
    setTimeout(() => { loadJobs(); }, 3000);
  } catch(e) {
    showAlert('Error: ' + e.message, false);
  } finally {
    btn.disabled = false;
    btn.innerHTML = '&#x25B6; Encolar 5 jobs de demo';
  }
}

async function enqueueAlert() {
  try {
    const r = await fetch(BASE + '/api/jobs', {
      method: 'POST',
      headers: {'Content-Type':'application/json'},
      body: JSON.stringify({
        type: 'send_alert',
        payload: {to:'ops@demo.com', message:'Prueba desde el portfolio de Leonardo Gómez', level:'warning'},
        priority: 8
      })
    });
    const d = await r.json();
    showAlert('&#x1F6A8; Alerta encolada con ID #' + d.job.id + ' — se procesar&aacute; en segundos', true);
    setTimeout(() => loadJobs(), 3000);
  } catch(e) {
    showAlert('Error: ' + e.message, false);
  }
}

function showAlert(msg, ok) {
  const el = document.getElementById('demo-alert');
  el.className = 'alert show ' + (ok ? 'alert-ok' : 'alert-err');
  el.innerHTML = msg;
  setTimeout(() => { el.classList.remove('show'); }, 5000);
}

// Cargar al inicio y cada 10 segundos
loadStats();
loadJobs();
setInterval(() => { loadStats(); loadJobs(); }, 10000);
</script>
</body>
</html>`
