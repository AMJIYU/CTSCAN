export interface LogColumn {
  key: string
  label: string
}

export interface LogKeyValueItem {
  label: string
  value: unknown
}

export interface LogSection {
  title: string
  columns?: LogColumn[]
  rows?: Array<Record<string, unknown>>
  items?: LogKeyValueItem[]
  content?: string
}

export interface LogSnapshot {
  title: string
  description?: string
  filters?: Record<string, unknown>
  sections: LogSection[]
}

export const createDefaultLogFileName = (title: string, extension = 'html') => {
  const safeTitle = sanitizeFileName(title || 'CTScan日志')
  const safeExtension = extension.replace(/^\.+/, '') || 'html'
  const timestamp = new Date()
    .toISOString()
    .replace(/\.\d{3}Z$/, '')
    .replace(/[:T]/g, '-')
  return `${safeTitle}_${timestamp}.${safeExtension}`
}

export const buildLogMarkdown = (snapshot: LogSnapshot) => {
  const exportedAt = formatDateTime(new Date())
  const lines: string[] = [
    `# ${snapshot.title || 'CTScan 日志'}`,
    '',
    `- 导出时间：${exportedAt}`,
    '- 导出工具：CTScan',
    ''
  ]

  if (snapshot.description) {
    lines.push(snapshot.description, '')
  }

  if (snapshot.filters && hasMeaningfulValues(snapshot.filters)) {
    lines.push('## 筛选条件', '')
    lines.push(...buildKeyValueTable(snapshot.filters))
    lines.push('')
  }

  snapshot.sections.forEach(section => {
    lines.push(`## ${section.title}`, '')

    if (section.items) {
      lines.push(...buildKeyValueItems(section.items))
      lines.push('')
    }

    if (section.content) {
      lines.push('```text', section.content.trim() || '暂无内容', '```', '')
    }

    if (section.columns && section.rows) {
      lines.push(`共 ${section.rows.length} 条记录。`, '')
      lines.push(...buildMarkdownTable(section.columns, section.rows))
      lines.push('')
    }
  })

  lines.push('---', '本文件为 CTScan 手动保存的页面查询快照。')
  return `${lines.join('\n')}\n`
}

export const buildPortableLogHtml = (snapshot: LogSnapshot) => {
  const exportedAt = formatDateTime(new Date())
  const reportData = {
    ...snapshot,
    exportedAt,
    generator: 'CTScan',
    formatVersion: 1
  }
  const json = JSON.stringify(reportData)
    .replace(/</g, '\\u003c')
    .replace(/>/g, '\\u003e')
    .replace(/&/g, '\\u0026')

  return `<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>${escapeHtml(snapshot.title || 'CTScan 日志')}</title>
  <style>
    :root {
      color-scheme: light;
      --bg: #f4f7fb;
      --panel: #ffffff;
      --text: #172033;
      --muted: #667085;
      --line: #d9e2ef;
      --primary: #1f6feb;
      --primary-soft: #eaf2ff;
      --danger: #d92d20;
      --shadow: 0 14px 40px rgba(16, 24, 40, .08);
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      background: radial-gradient(circle at top left, #ffffff 0, var(--bg) 42%, #eef4fb 100%);
      color: var(--text);
      font-family: "Microsoft YaHei", "PingFang SC", "Segoe UI", sans-serif;
      line-height: 1.5;
    }
    .page { max-width: 1440px; margin: 0 auto; padding: 32px 28px 48px; }
    .hero {
      display: flex;
      justify-content: space-between;
      gap: 24px;
      padding: 28px;
      border: 1px solid var(--line);
      border-radius: 22px;
      background: rgba(255, 255, 255, .9);
      box-shadow: var(--shadow);
    }
    .hero h1 { margin: 0 0 10px; font-size: 30px; letter-spacing: -.02em; }
    .hero p { margin: 0; color: var(--muted); }
    .meta { display: grid; gap: 8px; min-width: 260px; color: var(--muted); font-size: 13px; }
    .meta strong { color: var(--text); }
    .toolbar {
      display: flex;
      flex-wrap: wrap;
      gap: 12px;
      align-items: center;
      margin: 22px 0;
      padding: 16px;
      border: 1px solid var(--line);
      border-radius: 16px;
      background: rgba(255, 255, 255, .88);
    }
    input[type="search"] {
      min-width: min(460px, 100%);
      flex: 1;
      padding: 11px 13px;
      border: 1px solid var(--line);
      border-radius: 12px;
      font-size: 14px;
      outline: none;
      background: #fff;
    }
    input[type="search"]:focus { border-color: var(--primary); box-shadow: 0 0 0 3px var(--primary-soft); }
    button {
      border: 1px solid var(--line);
      background: #fff;
      color: var(--text);
      border-radius: 12px;
      padding: 10px 14px;
      cursor: pointer;
      font-size: 14px;
    }
    button.primary { border-color: var(--primary); background: var(--primary); color: #fff; }
    .section {
      margin-top: 18px;
      padding: 22px;
      border: 1px solid var(--line);
      border-radius: 18px;
      background: var(--panel);
      box-shadow: 0 8px 26px rgba(16, 24, 40, .05);
    }
    .section h2 { margin: 0 0 14px; font-size: 20px; }
    .section-head {
      display: flex;
      justify-content: space-between;
      align-items: center;
      gap: 12px;
      margin-bottom: 12px;
    }
    .count { color: var(--muted); font-size: 13px; white-space: nowrap; }
    .table-wrap { overflow: auto; border: 1px solid var(--line); border-radius: 14px; }
    table { width: 100%; border-collapse: collapse; min-width: 760px; }
    th, td {
      padding: 10px 12px;
      border-bottom: 1px solid #edf1f7;
      text-align: left;
      vertical-align: top;
      font-size: 13px;
    }
    th {
      position: sticky;
      top: 0;
      z-index: 1;
      background: #f8fafc;
      color: #344054;
      font-weight: 700;
    }
    tr:hover td { background: #fafcff; }
    td { max-width: 460px; word-break: break-word; }
    .kv-table { min-width: 0; }
    .kv-table th:first-child, .kv-table td:first-child { width: 180px; color: var(--muted); }
    .empty { padding: 20px; color: var(--muted); text-align: center; }
    .pager {
      display: flex;
      align-items: center;
      justify-content: flex-end;
      gap: 10px;
      margin-top: 12px;
      color: var(--muted);
      font-size: 13px;
    }
    .footer { margin-top: 24px; color: var(--muted); font-size: 12px; text-align: center; }
    pre {
      white-space: pre-wrap;
      word-break: break-word;
      background: #101828;
      color: #f9fafb;
      border-radius: 14px;
      padding: 16px;
      overflow: auto;
    }
    @media print {
      body { background: #fff; }
      .toolbar, .pager { display: none; }
      .page { max-width: none; padding: 0; }
      .hero, .section { box-shadow: none; break-inside: avoid; }
      th { position: static; }
    }
  </style>
</head>
<body>
  <div class="page">
    <header class="hero">
      <div>
        <h1 id="report-title"></h1>
        <p id="report-description"></p>
      </div>
      <div class="meta">
        <div>导出时间：<strong id="report-time"></strong></div>
        <div>导出工具：<strong>CTScan</strong></div>
        <div>格式：<strong>自包含 HTML，可离线打开</strong></div>
      </div>
    </header>

    <div class="toolbar">
      <input id="global-search" type="search" placeholder="在本报告内搜索关键字、IP、用户名、路径、事件ID...">
      <button class="primary" id="download-json">下载原始 JSON</button>
      <button id="print-report">打印/另存为 PDF</button>
    </div>

    <main id="report-root"></main>
    <div class="footer">本文件为 CTScan 手动保存的页面查询快照；数据已内嵌在 HTML 中，复制到其他电脑后可直接打开。</div>
  </div>
  <script id="ctscan-report-data" type="application/json">${json}</script>
  <script>
    (function () {
      const data = JSON.parse(document.getElementById('ctscan-report-data').textContent || '{}')
      const state = {}
      const pageSize = 100
      const root = document.getElementById('report-root')
      const globalSearch = document.getElementById('global-search')

      document.getElementById('report-title').textContent = data.title || 'CTScan 日志'
      document.getElementById('report-description').textContent = data.description || ''
      document.getElementById('report-time').textContent = data.exportedAt || '-'
      document.getElementById('print-report').addEventListener('click', function () { window.print() })
      document.getElementById('download-json').addEventListener('click', function () {
        const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json;charset=utf-8' })
        const url = URL.createObjectURL(blob)
        const link = document.createElement('a')
        link.href = url
        link.download = (data.title || 'ctscan-log') + '.json'
        link.click()
        URL.revokeObjectURL(url)
      })
      globalSearch.addEventListener('input', render)

      function text(value) {
        if (value === undefined || value === null || value === '') return '-'
        if (Array.isArray(value)) return value.map(text).join(', ')
        if (typeof value === 'object') return JSON.stringify(value)
        return String(value)
      }

      function el(tag, className, content) {
        const node = document.createElement(tag)
        if (className) node.className = className
        if (content !== undefined) node.textContent = content
        return node
      }

      function renderKeyValues(title, values) {
        const section = el('section', 'section')
        section.appendChild(el('h2', '', title))
        const rows = Object.entries(values || {}).filter(function (entry) {
          return text(entry[1]).trim() !== '-'
        })
        if (!rows.length) {
          section.appendChild(el('div', 'empty', '暂无数据'))
          return section
        }
        const wrap = el('div', 'table-wrap')
        const table = el('table', 'kv-table')
        table.innerHTML = '<thead><tr><th>字段</th><th>值</th></tr></thead>'
        const tbody = document.createElement('tbody')
        rows.forEach(function (entry) {
          const tr = document.createElement('tr')
          tr.appendChild(el('td', '', entry[0]))
          tr.appendChild(el('td', '', text(entry[1])))
          tbody.appendChild(tr)
        })
        table.appendChild(tbody)
        wrap.appendChild(table)
        section.appendChild(wrap)
        return section
      }

      function rowMatches(row, columns, query) {
        if (!query) return true
        return columns.some(function (column) {
          return text(row[column.key]).toLowerCase().includes(query)
        })
      }

      function renderTableSection(sectionData, index, query) {
        const columns = sectionData.columns || []
        const rows = sectionData.rows || []
        if (!state[index]) state[index] = { page: 1 }
        const matchedRows = rows.filter(function (row) { return rowMatches(row, columns, query) })
        const totalPages = Math.max(1, Math.ceil(matchedRows.length / pageSize))
        state[index].page = Math.min(state[index].page, totalPages)
        const start = (state[index].page - 1) * pageSize
        const visibleRows = matchedRows.slice(start, start + pageSize)

        const section = el('section', 'section')
        const head = el('div', 'section-head')
        head.appendChild(el('h2', '', sectionData.title || '数据'))
        head.appendChild(el('div', 'count', '共 ' + matchedRows.length + ' 条记录' + (rows.length !== matchedRows.length ? ' / 原始 ' + rows.length + ' 条' : '')))
        section.appendChild(head)

        if (!matchedRows.length) {
          section.appendChild(el('div', 'empty', '暂无匹配记录'))
          return section
        }

        const wrap = el('div', 'table-wrap')
        const table = document.createElement('table')
        const thead = document.createElement('thead')
        const headerRow = document.createElement('tr')
        columns.forEach(function (column) { headerRow.appendChild(el('th', '', column.label || column.key)) })
        thead.appendChild(headerRow)
        table.appendChild(thead)
        const tbody = document.createElement('tbody')
        visibleRows.forEach(function (row) {
          const tr = document.createElement('tr')
          columns.forEach(function (column) { tr.appendChild(el('td', '', text(row[column.key]))) })
          tbody.appendChild(tr)
        })
        table.appendChild(tbody)
        wrap.appendChild(table)
        section.appendChild(wrap)

        const pager = el('div', 'pager')
        const prev = el('button', '', '上一页')
        const next = el('button', '', '下一页')
        const info = el('span', '', '第 ' + state[index].page + ' / ' + totalPages + ' 页，每页 ' + pageSize + ' 条')
        prev.disabled = state[index].page <= 1
        next.disabled = state[index].page >= totalPages
        prev.addEventListener('click', function () { state[index].page -= 1; render() })
        next.addEventListener('click', function () { state[index].page += 1; render() })
        pager.appendChild(prev)
        pager.appendChild(info)
        pager.appendChild(next)
        section.appendChild(pager)
        return section
      }

      function renderSection(sectionData, index, query) {
        if (sectionData.columns && sectionData.rows) return renderTableSection(sectionData, index, query)
        const section = el('section', 'section')
        section.appendChild(el('h2', '', sectionData.title || '内容'))
        if (sectionData.items) {
          const table = renderKeyValues(sectionData.title || '字段', Object.fromEntries(sectionData.items.map(function (item) {
            return [item.label, item.value]
          })))
          return table
        }
        if (sectionData.content) {
          const pre = document.createElement('pre')
          pre.textContent = sectionData.content
          section.appendChild(pre)
        } else {
          section.appendChild(el('div', 'empty', '暂无数据'))
        }
        return section
      }

      function render() {
        const query = (globalSearch.value || '').trim().toLowerCase()
        root.innerHTML = ''
        if (data.filters && Object.values(data.filters).some(function (value) { return text(value).trim() !== '-' })) {
          root.appendChild(renderKeyValues('筛选条件', data.filters))
        }
        ;(data.sections || []).forEach(function (section, index) {
          root.appendChild(renderSection(section, index, query))
        })
      }

      render()
    })()
  </script>
</body>
</html>`
}

export const formatBytes = (value: number | undefined | null) => {
  if (!value) return value === 0 ? '0 B' : '-'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const index = Math.min(Math.floor(Math.log(value) / Math.log(1024)), units.length - 1)
  return `${(value / Math.pow(1024, index)).toFixed(2)} ${units[index]}`
}

const buildKeyValueTable = (values: Record<string, unknown>) => {
  const rows = Object.entries(values)
    .filter(([, value]) => isMeaningfulValue(value))
    .map(([key, value]) => `| ${escapeMarkdownCell(key)} | ${escapeMarkdownCell(formatValue(value))} |`)
  if (rows.length === 0) return ['暂无筛选条件。']
  return ['| 条件 | 值 |', '| --- | --- |', ...rows]
}

const buildKeyValueItems = (items: LogKeyValueItem[]) => {
  if (items.length === 0) return ['暂无数据。']
  return [
    '| 字段 | 值 |',
    '| --- | --- |',
    ...items.map(item => `| ${escapeMarkdownCell(item.label)} | ${escapeMarkdownCell(formatValue(item.value))} |`)
  ]
}

const buildMarkdownTable = (columns: LogColumn[], rows: Array<Record<string, unknown>>) => {
  if (rows.length === 0) return ['暂无记录。']
  const header = `| ${columns.map(column => escapeMarkdownCell(column.label)).join(' | ')} |`
  const divider = `| ${columns.map(() => '---').join(' | ')} |`
  const body = rows.map(row => {
    return `| ${columns.map(column => escapeMarkdownCell(formatValue(row[column.key]))).join(' | ')} |`
  })
  return [header, divider, ...body]
}

const formatValue = (value: unknown): string => {
  if (value === undefined || value === null || value === '') return '-'
  if (value instanceof Date) return formatDateTime(value)
  if (typeof value === 'boolean') return value ? '是' : '否'
  if (Array.isArray(value)) return value.map(formatValue).join(', ')
  if (typeof value === 'object') return JSON.stringify(value)
  return String(value)
}

const formatDateTime = (date: Date) => {
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false
  })
}

const escapeMarkdownCell = (value: string) => {
  return value.replace(/\|/g, '\\|').replace(/\r?\n/g, '<br>')
}

const escapeHtml = (value: string) => {
  return value
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;')
}

const hasMeaningfulValues = (values: Record<string, unknown>) => {
  return Object.values(values).some(isMeaningfulValue)
}

const isMeaningfulValue = (value: unknown) => {
  return value !== undefined && value !== null && String(value).trim() !== ''
}

const sanitizeFileName = (name: string) => {
  return name.replace(/[\\/:*?"<>|]/g, '-').replace(/\s+/g, '_').replace(/^[-_.]+|[-_.]+$/g, '')
}
