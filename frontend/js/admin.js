/* =====================================================
   JAMSEL COSMETICS — Admin Dashboard JavaScript
   Reads orders from localStorage, manages status & deletion
   ===================================================== */

// ============ STATE ============
let allOrders = [];
let currentFilter = 'all';

// ============ DOM READY ============
document.addEventListener('DOMContentLoaded', () => {
    loadOrders();
    initTabs();
    initMobileMenu();
    initOrderFilters();
    renderDashboard();
    renderOrdersTable();
});

// ============ LOAD ORDERS ============
function loadOrders() {
    allOrders = JSON.parse(localStorage.getItem('jamsel_orders')) || [];
}

function saveOrders() {
    localStorage.setItem('jamsel_orders', JSON.stringify(allOrders));
}

// ============ TAB NAVIGATION ============
function initTabs() {
    document.querySelectorAll('.sidebar-link[data-tab]').forEach(link => {
        link.addEventListener('click', (e) => {
            e.preventDefault();
            const tab = link.getAttribute('data-tab');

            // Update active link
            document.querySelectorAll('.sidebar-link[data-tab]').forEach(l => l.classList.remove('active'));
            link.classList.add('active');

            // Update active tab
            document.querySelectorAll('.tab-content').forEach(t => t.classList.remove('active'));
            document.getElementById(`tab-${tab}`).classList.add('active');

            // Update title
            document.getElementById('page-title').textContent =
                tab === 'dashboard' ? 'Dashboard' : 'All Orders';

            // Close mobile sidebar
            document.getElementById('sidebar').classList.remove('open');

            // Refresh data
            loadOrders();
            if (tab === 'dashboard') renderDashboard();
            if (tab === 'orders') renderOrdersTable();
        });
    });
}

// ============ MOBILE MENU ============
function initMobileMenu() {
    document.getElementById('menu-toggle').addEventListener('click', () => {
        document.getElementById('sidebar').classList.toggle('open');
    });
}

// ============ DASHBOARD ============
function renderDashboard() {
    // Stats
    const totalOrders = allOrders.length;
    const totalRevenue = allOrders.reduce((sum, o) => sum + (o.total || 0), 0);
    const delivered = allOrders.filter(o => o.status === 'delivered').length;
    const customers = new Set(allOrders.map(o => o.customer?.email)).size;

    document.getElementById('stat-total-orders').textContent = totalOrders;
    document.getElementById('stat-revenue').textContent = `Nu.${totalRevenue.toFixed(2)}`;
    document.getElementById('stat-delivered').textContent = delivered;
    document.getElementById('stat-customers').textContent = customers;

    // Recent orders table (last 5)
    const tbody = document.getElementById('dashboard-orders-body');
    const emptyEl = document.getElementById('dashboard-empty');

    if (allOrders.length === 0) {
        tbody.parentElement.parentElement.style.display = 'none';
        emptyEl.style.display = 'block';
        return;
    }

    tbody.parentElement.parentElement.style.display = 'block';
    emptyEl.style.display = 'none';

    const recent = [...allOrders].reverse().slice(0, 5);
    tbody.innerHTML = recent.map(order => createDashboardRow(order)).join('');
}

function createDashboardRow(order) {
    const productNames = (order.items || []).map(i => `${i.name} ×${i.qty}`).join(', ');
    const date = new Date(order.date).toLocaleDateString('en-US', {
        year: 'numeric', month: 'short', day: 'numeric'
    });

    return `
        <tr>
            <td><strong>${order.id}</strong></td>
            <td>${order.customer?.name || 'N/A'}</td>
            <td class="product-list-cell"><span>${productNames || 'N/A'}</span></td>
            <td><strong>Nu.${(order.total || 0).toFixed(2)}</strong></td>
            <td>${date}</td>
            <td><span class="status-badge ${order.status}">${order.status}</span></td>
            <td>
                <div class="action-btns">
                    <button class="action-btn view" onclick="viewOrder('${order.id}')">View</button>
                    ${getStatusActions(order)}
                </div>
            </td>
        </tr>
    `;
}

// ============ ORDERS TABLE ============
function initOrderFilters() {
    document.querySelectorAll('.order-filter-btn').forEach(btn => {
        btn.addEventListener('click', () => {
            document.querySelectorAll('.order-filter-btn').forEach(b => b.classList.remove('active'));
            btn.classList.add('active');
            currentFilter = btn.getAttribute('data-status');
            renderOrdersTable();
        });
    });
}

function renderOrdersTable() {
    const tbody = document.getElementById('all-orders-body');
    const emptyEl = document.getElementById('orders-empty');

    let filtered = currentFilter === 'all'
        ? allOrders
        : allOrders.filter(o => o.status === currentFilter);

    if (filtered.length === 0) {
        tbody.parentElement.parentElement.style.display = 'none';
        emptyEl.style.display = 'block';
        return;
    }

    tbody.parentElement.parentElement.style.display = 'block';
    emptyEl.style.display = 'none';

    const sorted = [...filtered].reverse();
    tbody.innerHTML = sorted.map(order => createFullOrderRow(order)).join('');
}

function createFullOrderRow(order) {
    const productNames = (order.items || []).map(i => `${i.name} ×${i.qty}`).join(', ');
    const date = new Date(order.date).toLocaleDateString('en-US', {
        year: 'numeric', month: 'short', day: 'numeric'
    });
    const address = `${order.customer?.address || ''}, ${order.customer?.city || ''}${order.customer?.zip ? ' ' + order.customer.zip : ''}`;
    const paymentLabel = { cod: 'Cash on Delivery', card: 'Card', wallet: 'Mobile Wallet' };

    return `
        <tr>
            <td><strong>${order.id}</strong></td>
            <td>${order.customer?.name || 'N/A'}</td>
            <td>${order.customer?.email || 'N/A'}</td>
            <td>${order.customer?.phone || 'N/A'}</td>
            <td style="max-width:150px;font-size:0.82rem;">${address}</td>
            <td class="product-list-cell"><span>${productNames || 'N/A'}</span></td>
            <td><strong>Nu.${(order.total || 0).toFixed(2)}</strong></td>
            <td>${paymentLabel[order.payment] || order.payment || 'N/A'}</td>
            <td>${date}</td>
            <td><span class="status-badge ${order.status}">${order.status}</span></td>
            <td>
                <div class="action-btns">
                    <button class="action-btn view" onclick="viewOrder('${order.id}')">View</button>
                    ${getStatusActions(order)}
                    <button class="action-btn delete" onclick="deleteOrder('${order.id}')">Delete</button>
                </div>
            </td>
        </tr>
    `;
}

function getStatusActions(order) {
    let actions = '';
    if (order.status === 'pending') {
        actions += `<button class="action-btn ship" onclick="updateStatus('${order.id}', 'shipped')">Ship</button>`;
    }
    if (order.status === 'shipped') {
        actions += `<button class="action-btn deliver" onclick="updateStatus('${order.id}', 'delivered')">Deliver</button>`;
    }
    if (order.status === 'delivered') {
        actions += `<button class="action-btn" style="background:#EDE8F5;color:#5A4BA0;cursor:default;">Done ✓</button>`;
    }
    return actions;
}

// ============ ORDER ACTIONS ============
function updateStatus(orderId, newStatus) {
    const order = allOrders.find(o => o.id === orderId);
    if (!order) return;
    order.status = newStatus;
    saveOrders();
    renderDashboard();
    renderOrdersTable();
    showAdminToast(`Order ${orderId} marked as ${newStatus}`, 'success');
}

function deleteOrder(orderId) {
    if (!confirm(`Are you sure you want to delete order ${orderId}?`)) return;
    allOrders = allOrders.filter(o => o.id !== orderId);
    saveOrders();
    renderDashboard();
    renderOrdersTable();
    showAdminToast(`Order ${orderId} deleted`, 'error');
}

// ============ ORDER DETAIL MODAL ============
function viewOrder(orderId) {
    const order = allOrders.find(o => o.id === orderId);
    if (!order) return;

    const date = new Date(order.date).toLocaleDateString('en-US', {
        year: 'numeric', month: 'long', day: 'numeric', hour: '2-digit', minute: '2-digit'
    });
    const paymentLabel = { cod: 'Cash on Delivery', card: 'Credit/Debit Card', wallet: 'Mobile Wallet' };

    const body = document.getElementById('order-modal-body');
    body.innerHTML = `
        <p class="modal-section-title">Order Info</p>
        <div class="modal-detail-row">
            <span class="modal-detail-label">Order ID</span>
            <span class="modal-detail-value"><strong>${order.id}</strong></span>
        </div>
        <div class="modal-detail-row">
            <span class="modal-detail-label">Date</span>
            <span class="modal-detail-value">${date}</span>
        </div>
        <div class="modal-detail-row">
            <span class="modal-detail-label">Status</span>
            <span class="modal-detail-value"><span class="status-badge ${order.status}">${order.status}</span></span>
        </div>
        <div class="modal-detail-row">
            <span class="modal-detail-label">Payment</span>
            <span class="modal-detail-value">${paymentLabel[order.payment] || order.payment || 'N/A'}</span>
        </div>

        <p class="modal-section-title">Customer Details</p>
        <div class="modal-detail-row">
            <span class="modal-detail-label">Name</span>
            <span class="modal-detail-value">${order.customer?.name || 'N/A'}</span>
        </div>
        <div class="modal-detail-row">
            <span class="modal-detail-label">Email</span>
            <span class="modal-detail-value">${order.customer?.email || 'N/A'}</span>
        </div>
        <div class="modal-detail-row">
            <span class="modal-detail-label">Phone</span>
            <span class="modal-detail-value">${order.customer?.phone || 'N/A'}</span>
        </div>
        <div class="modal-detail-row">
            <span class="modal-detail-label">Address</span>
            <span class="modal-detail-value">${order.customer?.address || ''}, ${order.customer?.city || ''} ${order.customer?.zip || ''}</span>
        </div>

        <p class="modal-section-title">Products Ordered</p>
        <div class="modal-products-list">
            ${(order.items || []).map(item => `
                <div class="modal-product-item">
                    <span>${item.name} × ${item.qty}</span>
                    <span>Nu.${(item.price * item.qty).toFixed(2)}</span>
                </div>
            `).join('')}
        </div>

        <p class="modal-section-title">Order Total</p>
        <div class="modal-detail-row">
            <span class="modal-detail-label">Subtotal</span>
            <span class="modal-detail-value">Nu.${(order.subtotal || 0).toFixed(2)}</span>
        </div>
        <div class="modal-detail-row">
            <span class="modal-detail-label">Shipping</span>
            <span class="modal-detail-value">${order.shipping === 0 ? 'Free' : '$' + (order.shipping || 0).toFixed(2)}</span>
        </div>
        <div class="modal-detail-row">
            <span class="modal-detail-label">Total</span>
            <span class="modal-detail-value" style="font-weight:700;font-size:1.1rem;color:#B5687E;">Nu.${(order.total || 0).toFixed(2)}</span>
        </div>
    `;

    document.getElementById('order-modal').classList.add('open');
}

function closeOrderModal() {
    document.getElementById('order-modal').classList.remove('open');
}

// Close modal on overlay click
document.getElementById('order-modal')?.addEventListener('click', (e) => {
    if (e.target === e.currentTarget) closeOrderModal();
});

// ESC to close
document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') closeOrderModal();
});

// ============ TOAST ============
function showAdminToast(message, type = '') {
    const container = document.getElementById('admin-toast-container');
    const toast = document.createElement('div');
    toast.className = `admin-toast ${type}`;
    toast.textContent = message;
    container.appendChild(toast);
    setTimeout(() => toast.remove(), 3000);
}
