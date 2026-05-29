package mailer

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

// Mailer envía emails via SMTP estándar de Go (sin dependencias externas)
type Mailer struct {
	host     string
	port     int
	user     string
	password string
	from     string
	enabled  bool
}

func New(host string, port int, user, password, from string) *Mailer {
	enabled := user != "" && password != ""
	if !enabled {
		log.Println("⚠️  SMTP no configurado — los emails se simularán en consola")
	}
	return &Mailer{host: host, port: port, user: user, password: password, from: from, enabled: enabled}
}

// send envía el email usando net/smtp estándar
func (m *Mailer) send(to, subject, body string) error {
	if !m.enabled {
		log.Printf("📧 [SIMULADO] Para: %s\n   Asunto: %s\n   Body: %s\n", to, subject, truncate(body, 120))
		return nil
	}

	auth := smtp.PlainAuth("", m.user, m.password, m.host)
	addr := fmt.Sprintf("%s:%d", m.host, m.port)

	msg := strings.Join([]string{
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=UTF-8",
		fmt.Sprintf("From: TaskQueue <%s>", m.from),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", subject),
		"",
		body,
	}, "\r\n")

	return smtp.SendMail(addr, auth, m.from, []string{to}, []byte(msg))
}

// SendEmail envía un email personalizado
func (m *Mailer) SendEmail(to, subject, body string) error {
	return m.send(to, subject, body)
}

// SendWelcome envía el email de bienvenida con template HTML
func (m *Mailer) SendWelcome(to, name string) error {
	subject := fmt.Sprintf("¡Bienvenido/a %s! 🎉", name)
	body := fmt.Sprintf(`
<div style="font-family:Arial,sans-serif;max-width:600px;margin:0 auto;padding:20px">
  <h1 style="color:#2563eb">¡Hola, %s! 👋</h1>
  <p>Tu cuenta fue creada exitosamente en <strong>TaskQueue</strong>.</p>
  <div style="background:#f0f9ff;border-left:4px solid #2563eb;padding:15px;margin:20px 0">
    <strong>¿Qué puedes hacer?</strong>
    <ul>
      <li>📧 Encolar emails para envío asíncrono</li>
      <li>📊 Generar reportes programados</li>
      <li>🚨 Recibir alertas críticas</li>
      <li>🔄 Reintentos automáticos con backoff exponencial</li>
    </ul>
  </div>
  <p style="color:#64748b;font-size:12px">Enviado por TaskQueue Worker</p>
</div>`, name)
	return m.send(to, subject, body)
}

// SendReport envía un reporte de estadísticas de la cola
func (m *Mailer) SendReport(to, period string, stats map[string]interface{}) error {
	subject := fmt.Sprintf("📊 Reporte %s — TaskQueue", period)

	rows := ""
	for status, data := range stats {
		if d, ok := data.(map[string]interface{}); ok {
			rows += fmt.Sprintf(
				"<tr><td style='padding:8px;border:1px solid #e2e8f0'>%s</td>"+
					"<td style='padding:8px;border:1px solid #e2e8f0;text-align:center'>%v</td>"+
					"<td style='padding:8px;border:1px solid #e2e8f0;text-align:center'>%vs</td></tr>",
				status, d["total"], d["avg_seconds"],
			)
		}
	}

	body := fmt.Sprintf(`
<div style="font-family:Arial,sans-serif;max-width:600px;margin:0 auto;padding:20px">
  <h2 style="color:#1e293b">📊 Reporte %s</h2>
  <table style="width:100%%;border-collapse:collapse;margin-top:15px">
    <thead>
      <tr style="background:#1e293b;color:white">
        <th style="padding:10px">Estado</th>
        <th style="padding:10px">Total</th>
        <th style="padding:10px">Tiempo promedio</th>
      </tr>
    </thead>
    <tbody>%s</tbody>
  </table>
</div>`, period, rows)

	return m.send(to, subject, body)
}

// SendAlert envía una alerta con nivel de severidad
func (m *Mailer) SendAlert(to, message, level string) error {
	colors := map[string]string{"info": "#3b82f6", "warning": "#f59e0b", "critical": "#ef4444"}
	emojis := map[string]string{"info": "ℹ️", "warning": "⚠️", "critical": "🚨"}

	color := colors[level]
	if color == "" {
		color = "#64748b"
	}
	emoji := emojis[level]
	if emoji == "" {
		emoji = "📢"
	}

	subject := fmt.Sprintf("%s Alerta [%s] — TaskQueue", emoji, level)
	body := fmt.Sprintf(`
<div style="font-family:Arial,sans-serif;max-width:600px;margin:0 auto;padding:20px">
  <div style="background:%s;color:white;padding:15px;border-radius:8px">
    <h2>%s Alerta %s</h2>
    <p style="font-size:18px">%s</p>
  </div>
</div>`, color, emoji, level, message)

	return m.send(to, subject, body)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
