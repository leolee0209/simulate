package render

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type RunData struct {
	Label   string
	History []ExportMetric
}

// ExportGraphHTML generates and saves an HTML file containing the graph data and simulation settings.
// The file is saved to the specified filePath (e.g., "output/simulation_results.html").
func ExportGraphHTML(filePath string, runs []RunData, settings SimulationSettings, settingsMarkdown string) error {
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	fi, err := os.Lstat(filePath)
	var perm os.FileMode = 0644
	if err == nil {
		perm = fi.Mode().Perm()
	}

	historyJSON := generateHistoryJSON(runs)

	htmlContent := generateHTMLContent(historyJSON, settings, settingsMarkdown)

	// Write to file
	if err := os.WriteFile(filePath, []byte(htmlContent), perm); err != nil {
		return err
	}

	return nil
}

func generateHistoryJSON(runs []RunData) string {
	if len(runs) == 0 {
		return "[]"
	}

	var sb strings.Builder
	sb.WriteString("[\n")
	for r, run := range runs {
		if r > 0 {
			sb.WriteString(",\n")
		}
		sb.WriteString(fmt.Sprintf("  {\"label\": \"%s\", \"data\": [\n", run.Label))
		for i, metric := range run.History {
			if i > 0 {
				sb.WriteString(",\n")
			}
			sb.WriteString(fmt.Sprintf("    {\"generation\": %d, \"avgSelfishHerdChance\": %.6f, \"survivorCount\": %d}",
				metric.Generation, metric.AvgSelfishHerdChance, metric.SurvivorCount))
		}
		sb.WriteString("\n  ]}")
	}
	sb.WriteString("\n]")
	return sb.String()
}

type SimulationSettings struct {
	WorldWidth        float64
	WorldHeight       float64
	InitialPreyCount  int
	InitialPredators  int
	PreySpeed         float64
	PredatorSpeed     float64
	PredatorRestTicks int
	PreyCatchRadius   float64
	GenerationCount   int
	TotalTicks        uint64
	ExportTime        string
}

func generateSettingsJSON(settings SimulationSettings) string {
	return fmt.Sprintf(`{
  "worldWidth": %.1f,
  "worldHeight": %.1f,
  "initialPreyCount": %d,
  "initialPredators": %d,
  "preySpeed": %.2f,
  "predatorSpeed": %.2f,
  "predatorRestTicks": %d,
  "preyCatchRadius": %.1f,
  "generationCount": %d,
  "totalTicks": %d,
  "exportTime": "%s"
}`,
		settings.WorldWidth,
		settings.WorldHeight,
		settings.InitialPreyCount,
		settings.InitialPredators,
		settings.PreySpeed,
		settings.PredatorSpeed,
		settings.PredatorRestTicks,
		settings.PreyCatchRadius,
		settings.GenerationCount,
		settings.TotalTicks,
		settings.ExportTime,
	)
}

func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}

func generateHTMLContent(historyJSON string, settings SimulationSettings, settingsMarkdown string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>Simulation Results</title>
<script src="https://cdn.jsdelivr.net/npm/chart.js@3.9.1/dist/chart.min.js"></script>
<style>
body { font-family: monospace; margin: 20px; line-height: 1.6; }
h1 { border-bottom: 1px solid #000; padding-bottom: 5px; }
h2 { border-bottom: 1px solid #999; margin-top: 30px; padding-bottom: 3px; }
table { border-collapse: collapse; margin: 20px 0; }
th, td { border: 1px solid #ccc; padding: 8px; text-align: left; }
th { background-color: #f0f0f0; }
pre { background-color: #f5f5f5; padding: 10px; overflow-x: auto; }
canvas { max-width: 100%%; margin: 20px 0; }
.settings-table { margin-left: 20px; }
</style>
</head>
<body>
<h1>Simulation Results</h1>
<p>Exported: <strong>%s</strong></p>

<h2>Run Settings</h2>
<table class="settings-table">
<tr><td>World Size</td><td>%.1f x %.1f</td></tr>
<tr><td>Initial Prey</td><td>%d</td></tr>
<tr><td>Initial Predators</td><td>%d</td></tr>
<tr><td>Prey Speed</td><td>%.2f</td></tr>
<tr><td>Predator Speed</td><td>%.2f</td></tr>
<tr><td>Predator Rest Ticks</td><td>%d</td></tr>
<tr><td>Generations Completed</td><td>%d</td></tr>
<tr><td>Total Ticks</td><td>%d</td></tr>
</table>

<h2>Generation Data Graph</h2>
<canvas id="chart" width="1000" height="300"></canvas>

<h2>Simulation Settings Documentation</h2>
<pre>%s</pre>

<h2>Generation Data Table</h2>
<table>
<thead><tr><th>Run</th><th>Generation</th><th>Avg SelfishHerdChance</th><th>Survivor Count</th></tr></thead>
<tbody id="dataTableBody"></tbody>
</table>

<script>
const historyData = %s;

// Set up colors for datasets
const colors = ['#e6194b', '#3cb44b', '#ffe119', '#4363d8', '#f58231', '#911eb4', '#46f0f0', '#f032e6', '#bcf60c', '#fabebe'];

// Populate data table
const tbody = document.getElementById('dataTableBody');
historyData.forEach((run, index) => {
  run.data.forEach(entry => {
    const row = document.createElement('tr');
    row.innerHTML = '<td>' + run.label + '</td><td>' + entry.generation + '</td><td>' + entry.avgSelfishHerdChance.toFixed(4) + '</td><td>' + entry.survivorCount + '</td>';
    tbody.appendChild(row);
  });
});

// Draw chart
const ctx = document.getElementById('chart').getContext('2d');
const datasets = historyData.map((run, index) => {
  const color = colors[index %% colors.length];
  return {
    label: run.label,
    data: run.data.map(d => d.avgSelfishHerdChance),
    borderColor: color,
    backgroundColor: color + '33', // 20%% opacity
    borderWidth: 1,
    fill: false,
    pointRadius: 2
  };
});

// Assuming all runs have the same generations, use the first run for labels
const labels = historyData.length > 0 ? historyData[0].data.map(d => d.generation) : [];

new Chart(ctx, {
  type: 'line',
  data: {
    labels: labels,
    datasets: datasets
  },
  options: {
    responsive: false,
    plugins: { legend: { display: true } },
    scales: {
      y: { min: -1, max: 1 },
      x: { title: { display: true, text: 'Generation' } }
    }
  }
});
</script>
</body>
</html>`,
		settings.ExportTime,
		settings.WorldWidth,
		settings.WorldHeight,
		settings.InitialPreyCount,
		settings.InitialPredators,
		settings.PreySpeed,
		settings.PredatorSpeed,
		settings.PredatorRestTicks,
		settings.GenerationCount,
		settings.TotalTicks,
		escapeHTML(settingsMarkdown),
		historyJSON,
	)
}
