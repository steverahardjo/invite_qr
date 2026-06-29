export interface LogEntry {
  id: number
  method: string
  path: string
  status: 'pending' | 'success' | 'error'
  timestamp: Date
  error?: string
}

let entries: LogEntry[] = []
let nextId = 0
type Listener = () => void
const listeners = new Set<Listener>()

export function logApiCall(
  method: string,
  path: string,
  status: LogEntry['status'],
  error?: string,
) {
  const entry: LogEntry = {
    id: nextId++,
    method,
    path,
    status,
    timestamp: new Date(),
    error,
  }
  entries = [...entries.slice(-49), entry]

  const icon = status === 'success' ? '✓' : status === 'error' ? '✗' : '…'
  const ts = entry.timestamp.toLocaleTimeString()
  const style =
    status === 'success'
      ? 'color:#9CAF88;font-weight:600'
      : status === 'error'
        ? 'color:#B76E79;font-weight:600'
        : 'color:#C9A96E;font-weight:600'
  console.log(
    `%c[API]%c ${icon} ${method} ${path}${error ? ' — ' + error : ''} %c${ts}`,
    'color:#C9A96E;font-weight:700',
    style,
    'color:#8B7D72;font-size:0.75em',
  )

  listeners.forEach((fn) => fn())
  return entry
}

export function subscribe(fn: Listener) {
  listeners.add(fn)
  return () => { listeners.delete(fn) }
}

export function getEntries(): LogEntry[] {
  return entries
}
