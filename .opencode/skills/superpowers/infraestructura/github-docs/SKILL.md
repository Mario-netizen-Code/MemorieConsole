---
name: github-docs
description: Se activa ÚNICAMENTE cuando el usuario menciona repositorios de GitHub, URLs de git o solicita investigar la instalación de librerías externas.
license: MIT
compatibility: opencode
---

## Condición de Activación

- El usuario usa palabras clave como: "github", "repositorio", "instalar la librería X", "url", "https://github.com...", o nombres en formato "usuario/proyecto" (ej. opendatalab/MinerU).
- Si la consulta es sobre código local, bugs existentes o lógica pura, **IGNORA esta skill** y responde directamente.

## Flujo de Ejecución

1. **Frenar el Contexto Principal:** No uses herramientas de red en este chat. Mantén el contexto limpio.
2. **Crear Subagente de Red:** Instancia un subagente investigador con:
   > "Eres un subagente extractor. Conéctate a la API de GitHub o descarga el README del repositorio. Busca cómo se instala (pip, conda, npm) y qué hace en 1 frase. Guarda el resultado en `.opencode/install_specs.json`."
3. **Leer archivo local:** Una vez que el subagente cree `.opencode/install_specs.json`, léelo.
4. **Responder al Usuario:** Usa los datos del JSON para explicar la instalación.

## Formato esperado (`.opencode/install_specs.json`)

```json
{
  "package_name": "Nombre",
  "summary": "Qué hace en una frase.",
  "install_commands": ["comando 1", "comando 2"]
}
```
