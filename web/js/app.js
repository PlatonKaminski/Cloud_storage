// web/js/app.js

const API_BASE = '/api/v1'

let authToken = localStorage.getItem('authToken')

const authModal = document.getElementById('authModal')
const mainContent = document.getElementById('mainContent')
const loginForm = document.getElementById('loginForm')
const registerForm = document.getElementById('registerForm')
const authError = document.getElementById('authError')
const fileInput = document.getElementById('fileInput')

document.addEventListener('DOMContentLoaded', () => {
    setupAuthTabs()
    setupAuthForms()

    if (authToken) {
        showMainContent()
        renderExplorer()
    } else {
        showAuthModal()
    }
})

function setupAuthTabs() {
    const tabs = document.querySelectorAll('.tab-btn')
    tabs.forEach(tab => {
        tab.addEventListener('click', () => {
            const targetTab = tab.dataset.tab
            tabs.forEach(t => t.classList.remove('active'))
            tab.classList.add('active')
            document.querySelectorAll('.auth-form').forEach(f => f.classList.remove('active'))
            if (targetTab === 'login') loginForm.classList.add('active')
            else registerForm.classList.add('active')
            hideAuthError()
        })
    })
}

function setupAuthForms() {
    loginForm.addEventListener('submit', handleLogin)
    registerForm.addEventListener('submit', handleRegister)
}

async function readError(res) {
    const text = await res.text()
    try {
        const json = JSON.parse(text)
        return json.error || json.message || JSON.stringify(json)
    } catch {
        return text || `Ошибка ${res.status}`
    }
}

async function handleLogin(e) {
    e.preventDefault()
    const email    = document.getElementById('loginEmail').value.trim()
    const password = document.getElementById('loginPassword').value
    const btn      = e.target.querySelector('button[type="submit"]')
    btn.disabled   = true
    const orig     = btn.textContent
    btn.textContent = 'Вход...'

    try {
        const res = await fetch('/auth/login', {
            method:  'POST',
            headers: { 'Content-Type': 'application/json' },
            body:    JSON.stringify({ email, password })
        })
        if (!res.ok) throw new Error(await readError(res))

        const data = await res.json()
        authToken  = data.token
        localStorage.setItem('authToken', authToken)
        hideAuthModal()
        showMainContent()
        renderExplorer()
        loginForm.reset()
        showMessage('✅ Успешный вход')
    } catch (err) {
        showAuthError(err.message)
    } finally {
        btn.disabled    = false
        btn.textContent = orig
    }
}

async function handleRegister(e) {
    e.preventDefault()
    const name     = document.getElementById('registerName').value.trim()
    const email    = document.getElementById('registerEmail').value.trim()
    const password = document.getElementById('registerPassword').value

    if (password.length < 6) { showAuthError('Минимум 6 символов'); return }

    const btn = e.target.querySelector('button[type="submit"]')
    btn.disabled = true
    const orig   = btn.textContent
    btn.textContent = 'Регистрация...'

    try {
        const res = await fetch('/auth/register', {
            method:  'POST',
            headers: { 'Content-Type': 'application/json' },
            body:    JSON.stringify({ name, email, password })
        })
        if (!res.ok) throw new Error(await readError(res))

        showMessage('✅ Регистрация успешна. Войдите в систему')
        document.querySelector('[data-tab="login"]').click()
        document.getElementById('loginEmail').value = email
        registerForm.reset()
    } catch (err) {
        showAuthError(err.message)
    } finally {
        btn.disabled    = false
        btn.textContent = orig
    }
}

function logout() {
    authToken = null
    localStorage.removeItem('authToken')
    loginForm.reset()
    registerForm.reset()
    hideMainContent()
    showAuthModal()
    showMessage('👋 Вы вышли из системы')
}


async function renderExplorer() {
    const explorer = document.getElementById('explorerView')
    explorer.innerHTML = '<p style="color:#999;padding:20px 0;">Загрузка...</p>'

    try {
        // Параллельно грузим папки и файлы
        const [foldersRes, filesRes] = await Promise.all([
            fetch(`${API_BASE}/folders`, { headers: authHeaders() }),
            fetch(`${API_BASE}/files`,   { headers: authHeaders() })
        ])

        if (foldersRes.status === 401 || filesRes.status === 401) { logout(); return }
        if (!foldersRes.ok) throw new Error(await readError(foldersRes))
        if (!filesRes.ok)   throw new Error(await readError(filesRes))

        const foldersData = await foldersRes.json()
        const filesData   = await filesRes.json()

// ✅ || [] защищает от null, undefined, false
        const folders = foldersData?.folders || foldersData?.data || (Array.isArray(foldersData) ? foldersData : [])
        const files   = filesData?.files    || filesData?.data   || (Array.isArray(filesData)   ? filesData   : [])

        // Обновляем select для загрузки
        updateFolderSelect(folders)

        // Строим HTML проводника
        explorer.innerHTML = buildExplorerHTML(folders, files)

    } catch (err) {
        explorer.innerHTML = `<p style="color:#dc3545;">❌ ${err.message}</p>`
    }
}

function updateFolderSelect(folders) {
    const sel = document.getElementById('folderSelect')
    sel.innerHTML = '<option value="">📁 Без папки (корень)</option>' +
        folders.map(f => `<option value="${esc(f.id)}">📁 ${esc(f.name)}</option>`).join('')
}

// ============================================================
// ПОСТРОЕНИЕ HTML ПРОВОДНИКА
// ============================================================

function buildExplorerHTML(folders, files) {
    const safeFolders = folders || []
    const safeFiles   = files   || []
    const rootFiles   = safeFiles.filter(f => !f.folder_id)

    let html = `
        <div class="explorer-header">
            <span>Название</span>
            <span style="text-align:right">Размер</span>
            <span>Тип</span>
            <span></span>
        </div>`

    safeFolders.forEach(folder => {
        const folderFiles = safeFiles.filter(f =>
            f.folder_id && String(f.folder_id) === String(folder.id)
        )
        html += renderFolder(folder, folderFiles)
    })

    rootFiles.forEach(file => {
        html += renderFile(file)
    })

    if (safeFolders.length === 0 && rootFiles.length === 0) {
        html = `
        <div class="empty-state">
            <div class="empty-state-icon">☁️</div>
            <div class="empty-state-text">Хранилище пусто</div>
            <div class="empty-state-sub">Загрузите первый файл или создайте папку</div>
        </div>`
    }

    return html
}

function renderFolder(folder, files) {
    const filesHTML = files.map(f => renderFile(f, true)).join('')

    return `
    <div class="explorer-folder">
        <div class="explorer-folder-header">
            <span class="folder-icon-name" onclick="toggleFolder('folder-${esc(folder.id)}')">
                <span class="folder-toggle" id="toggle-${esc(folder.id)}">▸</span>
                <span class="folder-emoji">📁</span>
                <span class="folder-title">${esc(folder.name)}</span>
                <span class="folder-count">${files.length} ${pluralFiles(files.length)}</span>
            </span>
            <div class="folder-meta"></div>
            <div class="folder-meta"></div>
            <div class="folder-actions">
                <button class="btn-icon" title="Удалить папку"
                    onclick="deleteFolder('${esc(folder.id)}', '${esc(folder.name)}')">
                    🗑️
                </button>
            </div>
        </div>
        <div class="explorer-folder-body" id="folder-${esc(folder.id)}" style="display:none;">
            ${filesHTML || `<div class="empty-folder">— Папка пуста</div>`}
        </div>
    </div>`
}

function renderFile(file, insideFolder = false) {
    return `
    <div class="explorer-file">
        <span class="file-icon-name">
            <span class="file-emoji">${getFileIcon(file.mime_type)}</span>
            <span class="file-name" title="${esc(file.name)}">${esc(file.name)}</span>
        </span>
        <span class="file-size">${formatFileSize(file.size)}</span>
        <span class="file-type">${getFileTypeName(file.mime_type)}</span>
        <div class="file-actions">
            <button class="btn-sm btn-primary"
                onclick="downloadFile('${esc(file.id)}', '${esc(file.name)}')">
                ⬇️ Скачать
            </button>
            <button class="btn-sm btn-danger"
                onclick="deleteFile('${esc(file.id)}', '${esc(file.name)}')">
                🗑️
            </button>
        </div>
    </div>`
}

function toggleFolder(id) {
    const body   = document.getElementById(id)
    const toggle = document.getElementById('toggle-' + id.replace('folder-', ''))
    const isOpen = body.style.display !== 'none'
    body.style.display  = isOpen ? 'none' : 'block'
    toggle.textContent  = isOpen ? '▸' : '▾'
}

// Добавьте эту функцию для красивых названий типов файлов
function getFileTypeName(mimeType) {
    if (!mimeType) return '—'
    if (mimeType.startsWith('image/')) return 'Изображение'
    if (mimeType.startsWith('video/')) return 'Видео'
    if (mimeType.startsWith('audio/')) return 'Аудио'
    if (mimeType.includes('pdf')) return 'PDF'
    if (mimeType.includes('word')) return 'Word'
    if (mimeType.includes('excel') || mimeType.includes('spreadsheet')) return 'Excel'
    if (mimeType.includes('zip') || mimeType.includes('rar') || mimeType.includes('tar')) return 'Архив'
    if (mimeType.startsWith('text/')) return 'Текст'
    return mimeType.split('/')[1]?.toUpperCase() || '—'
}
// ============================================================
// ФАЙЛЫ
// ============================================================

async function uploadFile() {
    const file     = fileInput.files[0]
    const folderId = document.getElementById('folderSelect').value

    if (!file) { showMessage('Выберите файл', true); return }

    // ✅ Смотрим что выбрано
    console.log('📁 Выбранная папка (folderId):', folderId)
    console.log('📁 Тип:', typeof folderId)

    const form = new FormData()
    form.append('file', file)

    // ✅ Проверяем что folderId не пустая строка
    if (folderId && folderId !== '') {
        form.append('folder_id', folderId)
        console.log('✅ folder_id добавлен в форму:', folderId)
    } else {
        console.log('⚠️ folder_id не добавлен — корень')
    }

    // ✅ Логируем всё что в FormData
    for (let [key, val] of form.entries()) {
        console.log(`FormData: ${key} =`, val)
    }

    try {
        const res = await fetch(`${API_BASE}/files`, {
            method:  'POST',
            headers: authHeaders(),
            body:    form
        })
        if (res.status === 401) { logout(); return }
        if (!res.ok) throw new Error(await readError(res))

        // ✅ Смотрим что вернул сервер
        const result = await res.json()
        console.log('📄 Ответ сервера после загрузки:', result)

        fileInput.value = ''
        showMessage(`✅ ${file.name} загружен`)
        await renderExplorer()
    } catch (err) {
        showMessage(`Ошибка: ${err.message}`, true)
    }
}

async function downloadFile(id, fileName) {
    try {
        const res = await fetch(`${API_BASE}/files/${id}`, { headers: authHeaders() })
        if (res.status === 401) { logout(); return }
        if (!res.ok) throw new Error(await readError(res))

        const blob = await res.blob()
        const url  = URL.createObjectURL(blob)
        const a    = document.createElement('a')
        a.href     = url
        a.download = fileName
        document.body.appendChild(a)
        a.click()
        document.body.removeChild(a)
        URL.revokeObjectURL(url)
        showMessage(`✅ ${fileName} скачан`)
    } catch (err) {
        showMessage(`Ошибка: ${err.message}`, true)
    }
}

async function deleteFile(id, fileName) {
    if (!confirm(`Удалить "${fileName}"?`)) return
    try {
        const res = await fetch(`${API_BASE}/files/${id}`, {
            method:  'DELETE',
            headers: authHeaders()
        })
        if (res.status === 401) { logout(); return }
        if (!res.ok) throw new Error(await readError(res))

        showMessage(`✅ ${fileName} удалён`)
        await renderExplorer()
    } catch (err) {
        showMessage(`Ошибка: ${err.message}`, true)
    }
}

// ============================================================
// ПАПКИ
// ============================================================

async function createFolder() {
    const input = document.getElementById('folderName')
    const name  = input.value.trim()
    if (!name) { showMessage('Введите название папки', true); return }

    try {
        const res = await fetch(`${API_BASE}/folders`, {
            method:  'POST',
            headers: { ...authHeaders(), 'Content-Type': 'application/json' },
            body:    JSON.stringify({ name })
        })
        if (res.status === 401) { logout(); return }
        if (!res.ok) throw new Error(await readError(res))

        input.value = ''
        showMessage(`✅ Папка "${name}" создана`)
        await renderExplorer()
    } catch (err) {
        showMessage(`Ошибка: ${err.message}`, true)
    }
}

async function deleteFolder(id, folderName) {
    if (!confirm(`Удалить папку "${folderName}" и все файлы в ней?`)) return
    try {
        const res = await fetch(`${API_BASE}/folders/${id}`, {
            method:  'DELETE',
            headers: authHeaders()
        })
        if (res.status === 401) { logout(); return }
        if (!res.ok) throw new Error(await readError(res))

        showMessage(`✅ Папка "${folderName}" удалена`)
        await renderExplorer()
    } catch (err) {
        showMessage(`Ошибка: ${err.message}`, true)
    }
}

// ============================================================
// УТИЛИТЫ
// ============================================================

function authHeaders() {
    return { 'Authorization': `Bearer ${authToken}` }
}

function esc(str) {
    if (!str) return ''
    return String(str)
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#39;')
}

function formatFileSize(bytes) {
    if (!bytes) return '0 B'
    const k = 1024, sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

function getFileIcon(mime) {
    if (!mime) return '📄'
    if (mime.startsWith('image/'))  return '🖼️'
    if (mime.startsWith('video/'))  return '🎬'
    if (mime.startsWith('audio/'))  return '🎵'
    if (mime.includes('pdf'))       return '📕'
    if (mime.includes('word'))      return '📘'
    if (mime.includes('excel') || mime.includes('spreadsheet')) return '📗'
    if (mime.includes('zip') || mime.includes('rar'))           return '📦'
    if (mime.startsWith('text/'))   return '📝'
    return '📄'
}

function pluralFiles(n) {
    if (n % 10 === 1 && n % 100 !== 11) return ''
    if ([2,3,4].includes(n % 10) && ![12,13,14].includes(n % 100)) return 'а'
    return 'ов'
}

function showMessage(text, isError = false) {
    document.querySelector('.message')?.remove()
    const msg = document.createElement('div')
    msg.className = 'message'
    msg.textContent = text
    msg.style.cssText = `
        position:fixed;bottom:20px;right:20px;padding:12px 20px;
        border-radius:8px;background:${isError ? '#dc3545' : '#28a745'};
        color:white;font-size:14px;z-index:10001;
        box-shadow:0 4px 12px rgba(0,0,0,0.15);`
    document.body.appendChild(msg)
    setTimeout(() => msg.remove(), 3000)
}

function showAuthModal()  { authModal.style.display = 'flex'; hideAuthError() }
function hideAuthModal()  { authModal.style.display = 'none' }
function showMainContent(){ mainContent.style.display = 'block' }
function hideMainContent(){ mainContent.style.display = 'none' }
function showAuthError(m) { authError.textContent = m; authError.classList.add('show') }
function hideAuthError()  { authError.textContent = ''; authError.classList.remove('show') }