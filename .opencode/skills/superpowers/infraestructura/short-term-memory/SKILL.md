---
name: short-term-memory
description: Use when you need to AUTOMATICALLY log every relevant step of a session in .memories/{project}_{date}.txt — implementations, modifications, errors, solutions, decisions, findings, verifications, and context — for later vector embedding and retrieval. This skill is ALWAYS active and runs automatically without being asked.
---

# Short-Term Memory

## ⚠️ AUTOMATIC BEHAVIOR — DO NOT SKIP

Esta skill se ejecuta **AUTOMÁTICAMENTE** durante toda la sesión. No esperes a que el usuario te pida registrar algo. Tan pronto como ocurra un evento relevante, ESCRIBELO inmediatamente.

**No es opcional. No preguntes. Simplemente hazlo.**

## ⚠️ FILE RULES — NEVER DELETE

- **NUNCA borres** archivos de memoria — perderías el historial para la vector DB
- **NUNCA reescribas** un archivo existente — solo haz **append**
- **NUNCA borres** entradas anteriores — cada entrada es un chunk valioso
- Si el archivo del día ya existe, haz append. Si no, créalo.
- Recomienda añadir `.memories/` a `.gitignore` (datos de sesión, no del proyecto)

## Directory Structure

```
.memories/
  convertfiletomarkdown_2026-05-27.txt
  convertfiletomarkdown_2026-05-28.txt
  otro-proyecto_2026-05-27.txt
```

Todos los archivos de memoria van dentro de `.memories/` en la raíz del proyecto.

## File Naming Convention

```
{project_name}_{YYYY-MM-DD}.txt
```

- `project_name` = nombre del proyecto. Detección automática según el tipo de proyecto:

  | Lenguaje | Archivo | Campo |
  |----------|---------|-------|
  | Python | `pyproject.toml` | `project.name` |
  | Python (legacy) | `setup.cfg` / `setup.py` | `metadata.name` / `name=` |
  | Node.js | `package.json` | `name` |
  | Rust | `Cargo.toml` | `package.name` |
  | Go | `go.mod` | `module` |
  | Java/Maven | `pom.xml` | `<artifactId>` |
  | Java/Gradle | `settings.gradle` | `rootProject.name` |
  | Ruby | `*.gemspec` | `spec.name` |
  | PHP | `composer.json` | `name` |
  | C/C++ | `CMakeLists.txt` | `project(...)` |
  | Elixir | `mix.exs` | `app:` |
  | .NET | `*.csproj` | `<PackageId>` o `<AssemblyName>` |
  | Genérico | nombre del directorio raíz | `basename` del repo |

  Busca los archivos en orden de prioridad. Usa el primero que encuentre.

- **Sanitización**: `project_name` se normaliza — minúsculas, espacios → guiones, sin acentos ni caracteres especiales. Ej: `"My Project!"` → `"my-project"`
- `YYYY-MM-DD` = fecha de la sesión (si la sesión cruza la medianoche, se crea un archivo nuevo)
- Ejemplo: `convertfiletomarkdown_2026-05-27.txt`

### File Rotation

- **Cada día = un archivo nuevo**. No acumules sesiones de diferentes días en el mismo archivo
- Al iniciar sesión, determina el `project_name` y la fecha actual, y escribe en `.memories/{project_name}_{YYYY-MM-DD}.txt`
- Si hoy es 2026-05-28, el archivo es `convertfiletomarkdown_2026-05-28.txt` — aunque ayer exista `...27.txt`
- El ID secuencial (`mem_001`) es **por archivo** — cada día arranca desde `mem_001`

## File Format

Cada entrada es UNA línea JSON minificada. Sin headings, sin separadores — auténtico JSONL:

```
{"id":"mem_001","timestamp":"2026-05-27T10:30:00+00:00","session_id":1,"scope":"project","important":false,"type":"modification","summary":"Fixed table extraction bug in formatter.py","files":["convertir_pdf/formatter.py"],"branch":"main","tags":["bugfix","pdf","tables"],"error":"IndexError: list index out of range in _extract_tables()","resolution":"Added colspan detection loop","details":"Found that merged cells in PDF tables caused IndexError. Added logic to handle colspan/rowspan by repeating cell values horizontally before vertical merge."}
```

Cada línea que empieza con `{` es una entrada completa. Líneas en blanco o que empiezan con `#` se ignoran durante el parseo.

### Field Rules

- **JSON minificado**: sin indentación, sin saltos de línea internos — **una sola línea por entrada**. `summary` debe ser conciso (máx 150 chars). `details` puede ser tan largo como sea necesario para capturar el contexto completo — no truncar información relevante
- `id` secuencial por archivo con **4 dígitos**: `mem_0001`, `mem_0002`, ..., `mem_0123`
- `session_id` entero que incrementa por sesión (para continuidad entre días). Ej: primera sesión = `1`, segunda = `2`...
- `scope` uno de: `agent` (memoria operativa del agente), `project` (decisión técnica del proyecto). Por defecto `project`
- `important` booleano — `true` si la entrada es crítica (cambios de schema, decisiones arquitectónicas). Por defecto `false`
- `timestamp` ISO 8601 **con timezone**: `YYYY-MM-DDTHH:MM:SS±HH:MM` o `...Z`
- `type` uno de: `context`, `implementation`, `modification`, `error`, `solution`, `decision`, `verification`, `finding`, `milestone`
- `summary` obligatorio — texto corto para retrieval semántico (máx 150 chars)
- `files` array de rutas relativas
- `branch` rama de git actual (`git branch --show-current`) — opcional pero recomendado
- `tags` array de strings para filtrado
- `error` solo si `type = error`, contiene el mensaje exacto
- `resolution` solo si aplica, describe cómo se solucionó
- `details` texto libre con el contexto completo (sin saltos de línea internos)
- No usar `\n` dentro del JSON — escapar como `\\n` si es necesario, o mejor usar `.` y `,`

## Memory Entry Types

| Type | Cuándo se escribe | Ejemplo de summary |
|------|-------------------|-------------------|
| `context` | Al inicio o cuando el usuario da requerimientos | "User wants daily file rotation with project name in .memories/" |
| `implementation` | Implementaste una feature o refactor completo | "Added --extract-all flag to CLI with tests" |
| `modification` | Cambiaste un archivo específico (edit/write) | "Fixed colspan handling in formatter.py table render" |
| `error` | Un comando falló (linter, typecheck, test, runtime) | "Ruff F401: pandas imported but unused" |
| `solution` | Resolviste un error o problema | "Removed unused pandas import to fix F401" |
| `decision` | Elegiste un enfoque entre varias opciones | "Chose JSONL over multi-line blocks for simpler parsing" |
| `verification` | Corriste linter/typecheck/tests y pasaron o fallaron | "Ruff check passed — 0 errors, 0 warnings" |
| `finding` | Descubriste algo durante exploración del código | "pdfplumber.PDF type from pdfplumber.pdf, not pdfplumber.PDF" |
| `milestone` | Completaste un hito o subobjetivo | "Migration from pdfminer to pdfplumber complete" |

### Rules per type

- **`:context`** → al recibir instrucciones del usuario o definir objetivos. Registra QUÉ se necesita hacer y POR QUÉ.
- **`:implementation`** → después de completar una feature o refactor. Incluye `files` con todos los archivos tocados.
- **`:modification`** → después de cada `edit` o `write`. Registrar qué cambiaste exactamente.
- **`:error`** → cuando un comando falla. Incluye el mensaje exacto en el campo `error`.
- **`:solution`** → después de arreglar un error. Relaciona con el `id` del error via `tags` o `details`.
- **`:decision`** → cuando eliges un approach. Incluye alternativas consideradas en `details`.
- **`:verification`** → después de ejecutar linter, typechecker o tests. Indica pass/fail.
- **`:finding`** → durante exploración. Descubrimientos sobre la arquitectura o restricciones.
- **`:milestone`** → cuando completas un objetivo intermedio o final.

No registres cada tool call individual. Registra cada **evento con significado**.

### Cómo escribir un buen `summary`

El summary es el campo más importante para retrieval semántico. Debe ser:

- **Específico**: no "Fixed a bug", sino "Fixed IndexError in formatter.py when table cells have colspan"
- **Con contexto**: incluir el módulo/área afectada ("extractor.py", "CLI", "tests")
- **Acción + objeto**: "Added X to Y", "Removed Z from W", "Changed A to B because C"

## Automatic Workflow

### Al iniciar sesión — auto-consulta + preparación

1. **Lee memoria reciente**: busca el archivo `.memories/{project_name}_*.txt` más reciente (por fecha en nombre). Lee las últimas **10 entradas** (líneas que empiezan con `{`). Esto te da continuidad con la sesión anterior sin necesidad de vector DB.
2. **Detecta `project_name`**: busca en orden los archivos de la tabla de detección — si nada coincide, usa el nombre del directorio raíz sanitizado.
3. **Detecta la rama actual**: `git branch --show-current`.
4. **Toma la fecha actual**: `YYYY-MM-DD`.
5. **Calcula `session_id`**: si hay archivos previos, `session_id = max_sesion_id_anterior + 1`. Si es la primera, `session_id = 1`.
6. **Construye la ruta**: `.memories/{project_name}_{YYYY-MM-DD}.txt`.
7. **Prepara el archivo**: si el directorio `.memories/` no existe, créalo. Si el archivo no existe, créalo. Recomienda `.gitignore` si no existe.
8. **Registra contexto**: escribe un entry `context` con el objetivo de la sesión si el usuario lo ha definido.

### Durante la sesión

9. **Después de cada acción relevante** (edit, write, run command, explore, fix error, etc.):
   - Evalúa si el resultado es significativo
   - Si lo es, genera el JSON minificado con todos los campos y haz append **INMEDIATAMENTE**
10. **Nunca preguntes** "¿Quieres que registre esto?" — simplemente hazlo
11. **Nunca acumules** entradas — cada entrada se escribe en el momento

### Al finalizar la sesión

12. **Registra fin de sesión**: cuando la sesión termine (el usuario se despide, cambia de tarea, o ejecuta un comando de cierre), escribe automáticamente un entry `milestone` con `summary: "Session ended — {resumen de lo logrado}"` e `important: true`. Esto marca un punto de corte claro para la próxima auto-consulta.

### Ejemplo de auto-consulta

```
Al iniciar sesión el 2026-05-28:
→ Lee .memories/convertfiletomarkdown_2026-05-27.txt
→ Últimas 10 entradas: mem_001 a mem_009
→ Encuentra: session_id=1, última entry "Session ended — skill improvements applied"
→ session_id=2, continúa desde ese punto
→ Escribe entry context con el nuevo objetivo
```

### Ejemplo de entrada con scope e important

```json
{"id":"mem_010","timestamp":"2026-05-28T09:00:00+00:00","session_id":2,"scope":"project","important":true,"type":"decision","summary":"Changed DB schema from SQLite to PostgreSQL","files":["migrations/"],"branch":"main","tags":["database","schema","breaking"],"details":"The vector DB needs concurrent write support. SQLite doesn't handle this well. Migrated to PostgreSQL with pgvector extension."}


### Minificar JSON (multi-lenguaje)

| Lenguaje | Código |
|----------|--------|
| Python | `json.dumps(obj, ensure_ascii=False, separators=(",",":"))` |
| JavaScript | `JSON.stringify(obj)` |
| Rust | `serde_json::to_string(&obj)?` |
| Go | `bytes, _ := json.Marshal(obj)` |
| Ruby | `obj.to_json` |
| PHP | `json_encode($obj, JSON_UNESCAPED_UNICODE)` |
| Java | `new Gson().toJson(obj)` o `ObjectMapper().writeValueAsString(obj)` |

### Ejemplo de flujo

```
Acción: usuario pide rotación diaria de archivos
→ .memories/convertfiletomarkdown_2026-05-27.txt:
  {"id":"mem_001","timestamp":"2026-05-27T15:00:00+00:00","session_id":1,"scope":"project","important":false,"type":"context","summary":"User wants daily file rotation with project name","files":[],"branch":"main","tags":["requirement","rotation"],"details":"..."}

Acción: editas formatter.py — cambio menor
→ Append al mismo archivo:
  {"id":"mem_002","timestamp":"2026-05-27T15:30:00+00:00","session_id":1,"scope":"project","important":false,"type":"modification","summary":"Fixed colspan in formatter.py table render","files":["convertir_pdf/formatter.py"],"branch":"main","tags":["bugfix","tables"],"details":"..."}

Acción: sesión termina
→ Append:
  {"id":"mem_010","timestamp":"2026-05-27T18:00:00+00:00","session_id":1,"scope":"project","important":true,"type":"milestone","summary":"Session ended — created short-term-memory skill with auto-review","files":["SKILL.md",".memories/"],"branch":"main","tags":["session-end","summary"],"details":"Completed: skill creation, format decisions, code review improvements, auto-consulta setup."}
```

## Integration with Vector DB

### Parseo (multi-lenguaje)

**Python:**
```python
import json
from pathlib import Path

def parse_all_memories(memories_dir: str = ".memories") -> list[dict]:
    entries = []
    for path in sorted(Path(memories_dir).glob("*.txt")):
        for line in path.read_text().splitlines():
            line = line.strip()
            if line.startswith("{"):
                entry = json.loads(line)
                entry["_source_file"] = path.name
                project, date = path.stem.rsplit("_", 1)
                entry["_project"] = project
                entry["_date"] = date
                entries.append(entry)
    return entries
```

**Node.js:**
```javascript
const fs = require("fs");
const path = require("path");
const memoriesDir = ".memories";
const files = fs.readdirSync(memoriesDir).filter(f => f.endsWith(".txt"));
const entries = [];
for (const file of files) {
  const lines = fs.readFileSync(path.join(memoriesDir, file), "utf8").split("\n");
  for (const line of lines) {
    if (line.startsWith("{")) entries.push({ ...JSON.parse(line), _source: file });
  }
}
```

### Embedding strategy

- **Texto a embed**: `f"{entry['summary']}\n{entry['details']}"`
- **Metadata del vector**: `{ "type": entry["type"], "scope": entry.get("scope","project"), "important": entry.get("important",False), "tags": entry["tags"], "files": entry["files"], "branch": entry.get("branch", ""), "timestamp": entry["timestamp"], "session_id": entry.get("session_id",0), "project": entry["_project"], "date": entry["_date"] }`
- **ID único vector DB**: `f"{entry['_project']}_{entry['_date']}_{entry['id']}"` — previene colisiones entre proyectos
- **Filtrado por proyecto**: usa el campo `_project` para limitar búsqueda al proyecto activo
- **Filtrado por importancia**: `where: {"important": True}` para recuperar solo decisiones críticas

## Quick Reference

| Acción | Formato |
|--------|---------|
| Directorio | `.memories/` en raíz del proyecto |
| Nombre archivo | `{project_name_sanitized}_{YYYY-MM-DD}.txt` |
| Sanitizar | minúsculas, espacios → `-`, sin acentos/símbolos |
| Detectar project | Busca config files del proyecto → fallback nombre directorio |
| Detectar branch | `git branch --show-current` |
| Calcular session_id | `max session_id de archivos previos + 1` |
| Auto-consulta inicio | Leer últimas 10 entradas del archivo más reciente |
| Formato entrada | Una línea JSON por entrada (JSONL) |
| Campos nuevos | `scope` (agent|project), `important` (bool), `session_id` (int) |
| Minificar JSON (Python) | `json.dumps(obj, ensure_ascii=False, separators=(",",":"))` |
| Minificar JSON (Node) | `JSON.stringify(obj)` |
| Parsear | `json.loads(line)` por cada línea que empiece con `{` |
| Timezone | Incluir en timestamp: `+00:00`, `-05:00`, o `Z` |
| Límite `summary` | 150 chars — `details` no tiene límite |
| ID vector DB | `{project}_{date}_{id}` |
| Fin de sesión | Entry `milestone` + `important:true` con resumen |
| `.gitignore` | Recomendar añadir `.memories/` |

## Common Mistakes

- **Solo registrar errores**: también registra context, implementaciones, verificaciones, hallazgos y milestones
- **JSON con saltos de línea**: usar `separators=(",",":")` en Python, `JSON.stringify()` en JS
- **Sin timezone**: timestamps sin timezone son ambiguos. Siempre incluir `+00:00` o `Z`
- **Summaries vagos**: "Fixed a bug" → "Fixed IndexError in extractor.py cuando celdas tienen colspan"
- **No registrar contexto**: el primer entry `context` con los requerimientos del usuario es el más valioso
- **No registrar branch**: sin branch no sabes en qué rama ocurrió cada cambio
- **No marcar entradas importantes**: decisiones arquitectónicas o cambios de schema deben tener `important: true`
- **Confundir scope agent vs project**: `scope: "agent"` es para memoria operativa del agente (por qué eligió un approach). `scope: "project"` es para decisiones técnicas del proyecto
- **No registrar decisiones de diseño**: son las memorias más valiosas para sesiones futuras
- **No registrar verificaciones exitosas**: "ruff check pasó" es tan valioso como "falló"
- **No hacer auto-consulta al inicio**: si no lees las últimas entradas, pierdes continuidad con la sesión anterior
- **No cerrar la sesión**: sin entry de fin de sesión, la próxima auto-consulta no sabe dónde terminó
- **Escribir en el archivo equivocado**: verificar `project_name`, fecha y branch antes de escribir
- **Mezclar días en un solo archivo**: cada día = su propio archivo
- **Detalles vacíos**: sin contexto semántico suficiente, el embedding es ruido
- **Registrar cada tool call**: solo eventos con valor semántico
- **Borrar o reescribir archivos**: siempre append, nunca reemplazar
- **No verificar project_name**: si cambias de proyecto, el archivo destino cambia
- **Summary demasiado largo**: >150 chars diluye el retrieval semántico. `details` puede ser extenso, el `summary` debe ser quirúrgico

## ⚠️ WINDOWS FILE LOCKING — ERROR AL ESCRIBIR

Cuando el **MemorieConsole tailer** está corriendo, los archivos de `.memories/` pueden quedar bloqueados para escritura por PowerShell.

### Síntoma

```
Add-Content : El proceso no puede obtener acceso al archivo
'...memorieconsole_2026-05-27.txt' porque está siendo utilizado en otro proceso.
```

### Causa

El tailer abre los archivos en modo lectura. En Windows, si el lector no abre con `FILE_SHARE_WRITE`, otros procesos (como `Add-Content`) no pueden escribir. El fix ya está aplicado en `file_windows.go` usando `syscall.CreateFile` con `FILE_SHARE_READ|FILE_SHARE_WRITE`.

Si aún así ocurre el error:

### Soluciones

| Prioridad | Solución | Cómo |
|-----------|----------|------|
| 1 | Usar .NET en vez de Add-Content | `[System.IO.File]::AppendAllText('ruta', "`n{...json...}")` |
| 2 | Esperar y reintentar | El handle se libera al cerrar el tailer (Ctrl+C). Esperar 2-3 segundos. |

### Automatización

Para escribir desde PowerShell mientras el tailer corre, usa siempre esta función:

```powershell
function Write-MemoryEntry {
    param($Path, $Json)
    [System.IO.File]::AppendAllText($Path, "`n$Json")
}
```

Luego: `Write-MemoryEntry -Path ".memories\proyecto_2026-05-27.txt" -Json '{...}'`
