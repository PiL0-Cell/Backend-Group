/* =====================================================
   JAMSEL COSMETICS — Shared JavaScript
   Backend API Integration with Session Cookies
   NO localStorage for authentication - all in database
   ===================================================== */

// ============ PRODUCT DATA (Fallback if backend fails) ============
const products = [
    {
        id: 1, name: "Velvet Rose Lipstick", category: "lipstick", price: 24.99, originalPrice: 32.99,
        image: "images/product_lipstick.png", rating: 4.8, reviews: 124, badge: "bestseller",
        isNew: false, isSale: true,
        description: "A luxuriously creamy matte lipstick in a universally flattering rose shade. Enriched with vitamin E and jojoba oil for long-lasting hydration and a velvety smooth finish that lasts up to 12 hours.",
        ingredients: "Ricinus Communis Seed Oil, Candelilla Wax, Isononyl Isononanoate, Vitamin E, Jojoba Oil, Iron Oxides, Rose Extract",
        shades: ["#C07070", "#B85C5C", "#D4899F", "#A0524D", "#8B4453"]
    },
    {
        id: 2, name: "Silk Glow Foundation", category: "foundation", price: 38.99, originalPrice: null,
        image: "images/product_foundation.png", rating: 4.9, reviews: 89, badge: "bestseller",
        isNew: false, isSale: false,
        description: "A lightweight, buildable foundation that melts into skin for a natural dewy finish. Infused with hyaluronic acid and SPF 30, it provides flawless coverage while keeping your skin hydrated all day long.",
        ingredients: "Water, Dimethicone, Glycerin, Hyaluronic Acid, Titanium Dioxide, Niacinamide, Zinc Oxide, Vitamin C",
        shades: ["#F5DEB3", "#DEB887", "#D2B48C", "#C8A882", "#A0845C"]
    },
    {
        id: 3, name: "Hydra Bloom Serum", category: "serum", price: 45.99, originalPrice: 58.99,
        image: "images/product_serum.png", rating: 4.9, reviews: 156, badge: "bestseller",
        isNew: false, isSale: true,
        description: "Our signature water-drop serum powered by triple hyaluronic acid complex. Delivers intense hydration deep into the skin, reducing fine lines and giving you a radiant, dewy glow from morning to night.",
        ingredients: "Aqua, Hyaluronic Acid (3 molecular weights), Niacinamide, Vitamin C, Rose Water, Aloe Vera, Glycerin, Squalane",
        shades: []
    },
    {
        id: 4, name: "Dew Drop Moisturizer", category: "moisturizer", price: 32.99, originalPrice: null,
        image: "images/product_moisturizer.png", rating: 4.7, reviews: 98, badge: null,
        isNew: false, isSale: false,
        description: "A cloud-soft gel-cream moisturizer that locks in 24-hour hydration. Formulated with ceramides and botanical extracts, it strengthens your skin barrier while leaving a silky, non-greasy finish.",
        ingredients: "Water, Ceramide Complex, Shea Butter, Squalane, Rosehip Oil, Aloe Vera Extract, Green Tea Extract, Vitamin E",
        shades: []
    },
    {
        id: 5, name: "Pure Petal Cleanser", category: "cleanser", price: 22.99, originalPrice: null,
        image: "images/product_cleanser.png", rating: 4.6, reviews: 73, badge: null,
        isNew: true, isSale: false,
        description: "A gentle foaming cleanser that removes makeup and impurities without stripping natural moisture. Infused with rose petal extract and chamomile for a soothing, spa-like cleansing experience.",
        ingredients: "Water, Cocamidopropyl Betaine, Rose Petal Extract, Chamomile Extract, Glycerin, Aloe Vera, Vitamin B5, Citric Acid",
        shades: []
    },
    {
        id: 6, name: "Solar Shield SPF 50", category: "sunscreen", price: 28.99, originalPrice: 35.99,
        image: "images/product_moisturizer.png", rating: 4.5, reviews: 67, badge: null,
        isNew: true, isSale: true,
        description: "A weightless, invisible sunscreen with SPF 50+ and PA++++. Blends seamlessly into all skin tones with zero white cast. Water-resistant formula that doubles as a hydrating primer.",
        ingredients: "Water, Homosalate, Octisalate, Zinc Oxide, Niacinamide, Hyaluronic Acid, Vitamin E, Centella Asiatica Extract",
        shades: []
    },
    {
        id: 7, name: "Starlight Eye Palette", category: "eye-makeup", price: 36.99, originalPrice: null,
        image: "images/product_lipstick.png", rating: 4.8, reviews: 112, badge: "bestseller",
        isNew: false, isSale: false,
        description: "A luxurious 12-shade eyeshadow palette featuring a curated mix of mattes, shimmers, and metallics in rose-gold, champagne, and bronze tones. Ultra-pigmented and blendable for everyday glam.",
        ingredients: "Mica, Talc, Dimethicone, Magnesium Stearate, Phenoxyethanol, Iron Oxides, Titanium Dioxide, Tin Oxide",
        shades: ["#E8C8B0", "#D4A084", "#C9A96E", "#B8860B", "#8B6F5C", "#A0524D"]
    },
    {
        id: 8, name: "Radiance Glow Set", category: "skincare-set", price: 89.99, originalPrice: 119.99,
        image: "images/products_flatlay.png", rating: 4.9, reviews: 201, badge: "bestseller",
        isNew: false, isSale: true,
        description: "Our best-selling skincare ritual in one beautiful gift set. Includes: Hydra Bloom Serum, Dew Drop Moisturizer, Pure Petal Cleanser, and a bonus Rose Water Mist. The complete glow routine.",
        ingredients: "See individual product pages for full ingredient lists.",
        shades: []
    },
    {
        id: 9, name: "Matte Luxe Lipstick", category: "lipstick", price: 21.99, originalPrice: null,
        image: "images/product_lipstick.png", rating: 4.6, reviews: 88, badge: null,
        isNew: true, isSale: false,
        description: "A transfer-proof matte lipstick with intense colour payoff. Fortified with argan oil to keep lips soft and comfortable. Available in 5 stunning nude and berry shades.",
        ingredients: "Isododecane, Dimethicone, Trimethylsiloxysilicate, Argan Oil, Vitamin E, Iron Oxides, Red 7 Lake",
        shades: ["#C48080", "#A0524D", "#8B4453", "#B5687E", "#D4899F"]
    },
    {
        id: 10, name: "Crystal Clear Toner", category: "cleanser", price: 26.99, originalPrice: null,
        image: "images/product_serum.png", rating: 4.5, reviews: 54, badge: null,
        isNew: true, isSale: false,
        description: "A gentle exfoliating toner with AHA and BHA to refine pores and brighten skin texture. pH-balanced formula with witch hazel and calendula for calm, clear, glowing skin.",
        ingredients: "Water, Glycolic Acid, Salicylic Acid, Witch Hazel, Calendula Extract, Panthenol, Allantoin, Citric Acid",
        shades: []
    },
    {
        id: 11, name: "Golden Hour Highlighter", category: "eye-makeup", price: 29.99, originalPrice: null,
        image: "images/product_foundation.png", rating: 4.7, reviews: 95, badge: null,
        isNew: false, isSale: false,
        description: "A baked powder highlighter in a warm champagne-gold shade. Finely milled for an ethereal, lit-from-within glow. Can be used on cheekbones, brow bone, and décolletage.",
        ingredients: "Mica, Talc, Synthetic Fluorphlogopite, Dimethicone, Iron Oxides, Tin Oxide, Titanium Dioxide",
        shades: ["#E8D5A8", "#C9A96E", "#F5DEB3"]
    },
    {
        id: 12, name: "Moisture Lock Cream", category: "moisturizer", price: 34.99, originalPrice: 42.99,
        image: "images/product_moisturizer.png", rating: 4.6, reviews: 76, badge: null,
        isNew: false, isSale: true,
        description: "A rich overnight cream that deeply nourishes and repairs skin while you sleep. Contains retinol, peptides, and squalane for firmer, plumper skin by morning.",
        ingredients: "Water, Squalane, Retinol, Peptide Complex, Shea Butter, Jojoba Oil, Vitamin E, Ceramide NP, Hyaluronic Acid",
        shades: []
    }
];

const sampleReviews = [
    { author: "Sonam P.", stars: 5, text: "Absolutely love this product! It's become a staple in my beauty routine. The quality is incredible.", date: "2 weeks ago" },
    { author: "Deki W.", stars: 5, text: "Beautiful packaging and even more beautiful results. My skin has never looked better!", date: "1 month ago" },
    { author: "Tshering L.", stars: 4, text: "Great product, very hydrating. I just wish they had more shade options. Will definitely repurchase.", date: "3 weeks ago" },
    { author: "Pema K.", stars: 5, text: "This is the best cosmetic product I've ever used. Worth every penny. Jamsel never disappoints!", date: "2 months ago" },
    { author: "Karma D.", stars: 4, text: "Lovely formula and lasts all day. The scent is so gentle and pleasant. Highly recommend.", date: "1 month ago" }
];

// ============ PATH DETECTION ============
const _inPages = window.location.pathname.includes('/pages/');
const pagePath = _inPages ? '' : 'pages/';
const imgPath = _inPages ? '../images/' : 'images/';

// ============ CART & WISHLIST STATE ============
let cart = [];
let wishlist = [];

// ============ HELPER: API CALL WITH CREDENTIALS ============
async function apiCall(url, method = 'GET', body = null) {
    const options = {
        method: method,
        credentials: 'include',  // Sends session cookie automatically
        headers: {
            'Content-Type': 'application/json'
        }
    };
    
    if (body) {
        options.body = JSON.stringify(body);
    }
    
    try {
        const response = await fetch(url, options);
        
        // Handle empty response
        const text = await response.text();
        if (!text) {
            return { ok: false, data: { error: 'Empty response from server' } };
        }
        
        const data = JSON.parse(text);
        return { ok: response.ok, data };
    } catch (error) {
        console.error('Fetch error:', error);
        return { ok: false, data: { error: 'Network error - is server running?' } };
    }
}

// ============ LOAD CART FROM BACKEND ============
async function loadCartFromBackend() {
    const result = await apiCall('/api/cart', 'GET');
    
    if (result.ok && Array.isArray(result.data)) {
        cart = result.data.map(item => ({
            id: item.product_id,
            qty: item.quantity
        }));
    } else {
        cart = [];
    }
}

// ============ LOAD WISHLIST FROM BACKEND ============
async function loadWishlistFromBackend() {
    const result = await apiCall('/api/wishlist', 'GET');
    
    if (result.ok && Array.isArray(result.data)) {
        wishlist = result.data.map(item => item.product_id);
    } else {
        wishlist = [];
    }
}

// ============ LOAD USER DATA FROM BACKEND ============
async function loadUserData() {
    const result = await apiCall('/api/user', 'GET');
    
    if (result.ok) {
        localStorage.setItem('jamsel_username', result.data.username);
        localStorage.setItem('jamsel_user_email', result.data.email);
        return true;
    } else {
        localStorage.removeItem('jamsel_username');
        localStorage.removeItem('jamsel_user_email');
        return false;
    }
}

// ============ CHECK IF USER IS LOGGED IN ============
async function isLoggedIn() {
    const result = await apiCall('/api/user', 'GET');
    return result.ok;
}

// ============ INIT (runs on every page) ============
document.addEventListener('DOMContentLoaded', async () => {
    initLoader();
    initStickyHeader();
    initSearch();
    initBackToTop();
    initModals();
    initFormListeners();
    
    // Load data from backend
    await loadUserData();
    await loadCartFromBackend();
    await loadWishlistFromBackend();
    
    updateCartBadge();
    updateWishlistBadge();
    updateAuthUI();
    const loggedIn = await isLoggedIn();
    if (loggedIn) {
        loadAIRecommendations();
    }
});

// ============ LOADING SCREEN ============
function initLoader() {
    const loader = document.getElementById('loading-screen');
    if (!loader) return;
    setTimeout(() => {
        loader.classList.add('hidden');
        setTimeout(() => loader.remove(), 600);
    }, 1500);
}

// ============ STICKY HEADER ============
function initStickyHeader() {
    const header = document.getElementById('main-header');
    if (!header) return;
    window.addEventListener('scroll', () => {
        header.classList.toggle('scrolled', window.scrollY > 50);
    });

    const hamburger = document.getElementById('hamburger-btn');
    const navLinks = document.getElementById('nav-links');
    if (hamburger && navLinks) {
        hamburger.addEventListener('click', () => {
            hamburger.classList.toggle('active');
            navLinks.classList.toggle('open');
        });
    }
}

// ============ SEARCH ============
function initSearch() {
    const toggleBtn = document.getElementById('search-toggle-btn');
    const searchBar = document.getElementById('search-bar');
    const closeBtn = document.getElementById('search-close-btn');
    const input = document.getElementById('search-input');
    const resultsContainer = document.getElementById('search-results');
    if (!toggleBtn || !searchBar) return;

    toggleBtn.addEventListener('click', () => {
        searchBar.classList.toggle('open');
        if (searchBar.classList.contains('open')) setTimeout(() => input.focus(), 300);
    });
    closeBtn.addEventListener('click', () => {
        searchBar.classList.remove('open');
        input.value = '';
        resultsContainer.innerHTML = '';
    });
    input.addEventListener('input', () => {
        const q = input.value.toLowerCase().trim();
        if (q.length < 2) { resultsContainer.innerHTML = ''; return; }
        const results = products.filter(p =>
            p.name.toLowerCase().includes(q) || p.category.toLowerCase().includes(q) || p.description.toLowerCase().includes(q)
        );
        if (results.length === 0) {
            resultsContainer.innerHTML = '<p style="padding:1rem;color:#999;text-align:center;">No products found</p>';
            return;
        }
        resultsContainer.innerHTML = results.map(p => `
            <div class="search-result-item" onclick="window.location.href='${pagePath}product-detail.html?id=${p.id}'">
                <img src="${imgPath}${p.image.replace('images/','')}" alt="${p.name}">
                <div class="sr-info"><h4>${p.name}</h4><p>Nu.${p.price.toFixed(2)}</p></div>
            </div>
        `).join('');
    });
}

// ============ BACK TO TOP ============
function initBackToTop() {
    const btn = document.getElementById('back-to-top');
    if (!btn) return;
    window.addEventListener('scroll', () => btn.classList.toggle('visible', window.scrollY > 600));
    btn.addEventListener('click', () => window.scrollTo({ top: 0, behavior: 'smooth' }));
}

// ============ MODALS ============
function initModals() {
    // Quick View Modal
    const qvModal = document.getElementById('quick-view-modal');
    const qvClose = document.getElementById('modal-close');
    if (qvModal && qvClose) {
        qvClose.addEventListener('click', closeQuickView);
        qvModal.addEventListener('click', e => { if (e.target === e.currentTarget) closeQuickView(); });
    }
    
    // ========== LOGIN/REGISTER MODAL ==========
    const loginBtn = document.getElementById('login-btn');
    const loginModal = document.getElementById('login-modal');
    const loginClose = document.getElementById('login-modal-close');
    
    if (loginBtn && loginModal) {
        loginBtn.addEventListener('click', async () => {
            // Check if already logged in via backend
            const loggedIn = await isLoggedIn();
            if (loggedIn) {
                window.location.href = '../../pages/account.html';
            } else {
                showLoginForm();
                loginModal.classList.add('open');
            }
        });
        
        loginClose.addEventListener('click', () => loginModal.classList.remove('open'));
        loginModal.addEventListener('click', e => { 
            if (e.target === e.currentTarget) loginModal.classList.remove('open'); 
        });
    }
    
    document.addEventListener('keydown', e => {
        if (e.key === 'Escape') {
            closeQuickView();
            if (loginModal) loginModal.classList.remove('open');
        }
    });
}

// ========== LOGIN FORM HTML ==========
function showLoginForm() {
    const modalContent = document.querySelector('#login-modal .login-content');
    if (!modalContent) return;
    
    modalContent.innerHTML = `
        <h2>Welcome Back</h2>
        <p>Sign in to your Jamsel account</p>
        <form class="login-form" id="login-form">
            <div class="form-group">
                <label for="login-email">Email Address *</label>
                <input type="email" id="login-email" required placeholder="your@email.com">
            </div>
            <div class="form-group">
                <label for="login-password">Password *</label>
                <input type="password" id="login-password" required placeholder="••••••••">
            </div>
            <button type="submit" class="btn-primary btn-full">Sign In</button>
        </form>
        <p class="login-note" style="margin-top:1rem; text-align:center;">
            Don't have an account? 
            <a href="#" id="show-register-link" style="color:#D4899F; text-decoration:underline; cursor:pointer;">Register here</a>
        </p>
    `;
    
    const loginForm = document.getElementById('login-form');
    if (loginForm) {
        loginForm.addEventListener('submit', handleLogin);
    }
    
    const registerLink = document.getElementById('show-register-link');
    if (registerLink) {
        registerLink.addEventListener('click', (e) => {
            e.preventDefault();
            showRegisterForm();
        });
    }
}

// ========== REGISTER FORM HTML ==========
function showRegisterForm() {
    const modalContent = document.querySelector('#login-modal .login-content');
    if (!modalContent) return;
    
    modalContent.innerHTML = `
        <h2>Create Account</h2>
        <p>Join the Jamsel family</p>
        <form class="login-form" id="register-form">
            <div class="form-group">
                <label for="reg-username">Username *</label>
                <input type="text" id="reg-username" required placeholder="Choose a username">
            </div>
            <div class="form-group">
                <label for="reg-email">Email Address *</label>
                <input type="email" id="reg-email" required placeholder="your@email.com">
            </div>
            <div class="form-group">
                <label for="reg-password">Password *</label>
                <input type="password" id="reg-password" required placeholder="••••••••">
            </div>
            <div class="form-group">
                <label for="reg-confirm-password">Confirm Password *</label>
                <input type="password" id="reg-confirm-password" required placeholder="••••••••">
            </div>
            <button type="submit" class="btn-primary btn-full">Create Account</button>
        </form>
        <p class="login-note" style="margin-top:1rem; text-align:center;">
            Already have an account? 
            <a href="#" id="show-login-link" style="color:#D4899F; text-decoration:underline; cursor:pointer;">Sign in here</a>
        </p>
    `;
    
    const registerForm = document.getElementById('register-form');
    if (registerForm) {
        registerForm.addEventListener('submit', handleRegister);
    }
    
    const loginLink = document.getElementById('show-login-link');
    if (loginLink) {
        loginLink.addEventListener('click', (e) => {
            e.preventDefault();
            showLoginForm();
        });
    }
}

// ========== HANDLE LOGIN ==========
async function handleLogin(e) {
    e.preventDefault();
    
    const email = document.getElementById('login-email').value.trim();
    const password = document.getElementById('login-password').value;
    
    if (!email || !password) {
        showToast('Please fill in all fields', 'error');
        return;
    }
    
    const result = await apiCall('/api/login', 'POST', { email, password });
    
    if (result.ok) {
        showToast(`Welcome back ${result.data.username}!`, 'success');
        document.getElementById('login-modal').classList.remove('open');
        
        await loadUserData();
        await loadCartFromBackend();
        await loadWishlistFromBackend();
        updateCartBadge();
        updateWishlistBadge();
        updateAuthUI();
        
        window.location.href = '../../pages/account.html';
    } else {
        showToast(result.data.error || 'Login failed', 'error');
    }
}

// ========== HANDLE REGISTER ==========
async function handleRegister(e) {
    e.preventDefault();
    
    const username = document.getElementById('reg-username').value.trim();
    const email = document.getElementById('reg-email').value.trim();
    const password = document.getElementById('reg-password').value;
    const confirmPassword = document.getElementById('reg-confirm-password').value;
    
    if (!username || !email || !password || !confirmPassword) {
        showToast('Please fill in all fields', 'error');
        return;
    }
    
    if (username.length < 3) {
        showToast('Username must be at least 3 characters', 'error');
        return;
    }
    
    if (!email.includes('@') || !email.includes('.')) {
        showToast('Please enter a valid email address', 'error');
        return;
    }
    
    if (password.length < 6) {
        showToast('Password must be at least 6 characters', 'error');
        return;
    }
    
    if (password !== confirmPassword) {
        showToast('Passwords do not match', 'error');
        return;
    }
    
    const result = await apiCall('/api/register', 'POST', { username, email, password });
    
    if (result.ok) {
        showToast(`Account created! Welcome ${username}!`, 'success');
        
        // Auto login after registration
        const loginResult = await apiCall('/api/login', 'POST', { email, password });
        
        if (loginResult.ok) {
            document.getElementById('login-modal').classList.remove('open');
            await loadUserData();
            updateAuthUI();
            window.location.href = '../../pages/account.html';
        }
    } else {
        showToast(result.data.error || 'Registration failed', 'error');
    }
}

// ========== UPDATE UI BASED ON LOGIN STATE ==========
async function updateAuthUI() {
    const loggedIn = await isLoggedIn();
    const loginBtn = document.getElementById('login-btn');
    const logoutBtn = document.getElementById('logout-btn');
    
    if (loginBtn) {
        if (loggedIn) {
            const username = localStorage.getItem('jamsel_username') || 'Account';
            loginBtn.innerHTML = `
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                    <circle cx="12" cy="7" r="4"/>
                </svg>
                <span style="font-size:0.75rem; margin-left:4px;">${username}</span>
            `;
            if (logoutBtn) logoutBtn.style.display = 'inline-flex';
        } else {
            loginBtn.innerHTML = `
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                    <circle cx="12" cy="7" r="4"/>
                </svg>
            `;
            if (logoutBtn) logoutBtn.style.display = 'none';
        }
    }
}

// ========== LOGOUT FUNCTION ==========
async function logout() {
    const result = await apiCall('/api/logout', 'POST');
    
    if (result.ok) {
        localStorage.removeItem('jamsel_username');
        localStorage.removeItem('jamsel_user_email');
        
        cart = [];
        wishlist = [];
        updateCartBadge();
        updateWishlistBadge();
        updateAuthUI();
        
        showToast('Logged out successfully!', 'success');
        window.location.reload();
    }
}

// ========== FORM LISTENERS ==========
function initFormListeners() {
    const nf = document.getElementById('newsletter-form');
    if (nf) nf.addEventListener('submit', e => { e.preventDefault(); showToast('Thank you for subscribing! 💌', 'success'); e.target.reset(); });
    
    const fnf = document.getElementById('footer-newsletter-form');
    if (fnf) fnf.addEventListener('submit', e => { e.preventDefault(); showToast('Thank you for subscribing! 💌', 'success'); e.target.reset(); });
    
    const cf = document.getElementById('contact-form');
    if (cf) cf.addEventListener('submit', e => { e.preventDefault(); showToast("Message sent! We'll get back to you soon. 📧", 'success'); e.target.reset(); });
    
    // Logout button
    const logoutBtn = document.getElementById('logout-btn');
    if (logoutBtn) {
        logoutBtn.addEventListener('click', logout);
    }
}

// ========== PRODUCT CARD RENDERING ==========
function renderProductGrid(containerId, items) {
    const container = document.getElementById(containerId);
    if (!container) return;
    container.innerHTML = items.map(p => createProductCard(p)).join('');
}

function createProductCard(product) {
    const isWishlisted = wishlist.includes(product.id);
    let badgeHTML = '';
    if (product.isSale && product.originalPrice) {
        const discount = Math.round(((product.originalPrice - product.price) / product.originalPrice) * 100);
        badgeHTML = `<span class="product-badge sale">${discount}% Off</span>`;
    } else if (product.isNew) {
        badgeHTML = `<span class="product-badge new">New</span>`;
    } else if (product.badge === 'bestseller') {
        badgeHTML = `<span class="product-badge bestseller">Best Seller</span>`;
    }
    const stars = '★'.repeat(Math.floor(product.rating)) + (product.rating % 1 >= 0.5 ? '½' : '');
    return `
        <div class="product-card" data-id="${product.id}">
            <div class="product-card-image">
                ${badgeHTML}
                <img src="${imgPath}${product.image.replace('images/','')}" alt="${product.name}" loading="lazy">
                <div class="product-card-actions">
                    <button class="product-action-btn ${isWishlisted ? 'wishlisted' : ''}" onclick="toggleWishlist(${product.id})" title="Wishlist">${isWishlisted ? '♥' : '♡'}</button>
                    <button class="product-action-btn" onclick="openQuickView(${product.id})" title="Quick View">👁</button>
                </div>
            </div>
            <div class="product-card-info">
                <p class="product-card-category">${product.category.replace('-', ' ')}</p>
                <a href="${pagePath}product-detail.html?id=${product.id}" class="product-card-name">${product.name}</a>
                <p class="product-card-rating">${stars} <span>(${product.reviews})</span></p>
                <div class="product-card-price">
                    <span class="current-price">Nu.${product.price.toFixed(2)}</span>
                    ${product.originalPrice ? `<span class="old-price">Nu.${product.originalPrice.toFixed(2)}</span>` : ''}
                </div>
                <button class="product-card-btn" onclick="addToCart(${product.id})">Add to Cart</button>
            </div>
        </div>
    `;
}

// ========== QUICK VIEW ==========
function openQuickView(id) {
    const p = products.find(pr => pr.id === id);
    if (!p) return;
    const content = document.getElementById('quick-view-content');
    if (!content) return;
    content.innerHTML = `
        <div class="quick-view-grid">
            <div class="qv-image"><img src="${imgPath}${p.image.replace('images/','')}" alt="${p.name}"></div>
            <div class="qv-info">
                <h2>${p.name}</h2>
                <p class="qv-price">Nu.${p.price.toFixed(2)} ${p.originalPrice ? `<span style="text-decoration:line-through;color:#999;font-size:0.9rem;">Nu.${p.originalPrice.toFixed(2)}</span>` : ''}</p>
                <p class="qv-desc">${p.description}</p>
                <button class="btn-primary" onclick="addToCart(${p.id}); closeQuickView();">Add to Cart</button>
                <a href="${pagePath}product-detail.html?id=${p.id}" class="btn-outline" style="margin-top:0.5rem;display:inline-block;text-align:center;width:100%;">View Details</a>
            </div>
        </div>
    `;
    document.getElementById('quick-view-modal').classList.add('open');
}

function closeQuickView() {
    const m = document.getElementById('quick-view-modal');
    if (m) m.classList.remove('open');
}

// ========== CART FUNCTIONS ==========
async function addToCart(id) {
    const loggedIn = await isLoggedIn();
    
    if (!loggedIn) {
        showToast('Please login first', 'error');
        document.getElementById('login-modal').classList.add('open');
        showLoginForm();
        return;
    }
    
    const result = await apiCall('/api/cart', 'POST', { product_id: id, quantity: 1 });
    
    if (result.ok) {
        await loadCartFromBackend();
        updateCartBadge();
        showToast('Added to cart!', 'success');
    } else {
        showToast(result.data.error || 'Failed to add to cart', 'error');
    }
}

async function removeFromCart(id) {
    const result = await apiCall(`/api/cart/${id}`, 'DELETE');
    
    if (result.ok) {
        await loadCartFromBackend();
        updateCartBadge();
        if (typeof renderCart === 'function') renderCart();
        showToast('Removed from cart');
    }
}

async function updateCartQty(id, delta) {
    const item = cart.find(i => i.id === id);
    if (!item) return;
    
    const newQty = item.qty + delta;
    
    if (newQty <= 0) {
        await removeFromCart(id);
        return;
    }
    
    const result = await apiCall(`/api/cart/${id}`, 'PUT', { quantity: newQty });
    
    if (result.ok) {
        await loadCartFromBackend();
        updateCartBadge();
        if (typeof renderCart === 'function') renderCart();
    }
}

function updateCartBadge() {
    const total = cart.reduce((sum, item) => sum + item.qty, 0);
    const badge = document.getElementById('cart-count');
    if (badge) { badge.textContent = total; badge.classList.toggle('show', total > 0); }
}

function getCartTotal() {
    return cart.reduce((sum, item) => {
        const p = products.find(pr => pr.id === item.id);
        return sum + (p ? p.price * item.qty : 0);
    }, 0);
}

// ========== WISHLIST FUNCTIONS ==========
async function removeFromWishlist(productId) {
    const result = await apiCall(`/api/wishlist/${productId}`, 'DELETE');
    
    if (result.ok) {
        wishlist = wishlist.filter(id => id !== productId);
        updateWishlistBadge();
        return true;
    }
    return false;
}

async function toggleWishlist(id) {
    const loggedIn = await isLoggedIn();
    
    if (!loggedIn) {
        showToast('Please login first', 'error');
        document.getElementById('login-modal').classList.add('open');
        showLoginForm();
        return;
    }
    
    const isInWishlist = wishlist.includes(id);
    
    let result;
    if (isInWishlist) {
        result = await apiCall(`/api/wishlist/${id}`, 'DELETE');
        if (result.ok) {
            wishlist = wishlist.filter(i => i !== id);
            showToast('Removed from wishlist');
        }
    } else {
        result = await apiCall('/api/wishlist', 'POST', { product_id: id });
        if (result.ok) {
            wishlist.push(id);
            showToast('Added to wishlist!', 'success');
        }
    }
    
    if (result && result.ok) {
        updateWishlistBadge();
        
        // Refresh product grids
        document.querySelectorAll('.products-grid').forEach(grid => {
            if (grid.children.length > 0) {
                const ids = Array.from(grid.querySelectorAll('.product-card')).map(c => parseInt(c.dataset.id));
                const items = products.filter(p => ids.includes(p.id));
                grid.innerHTML = items.map(p => createProductCard(p)).join('');
            }
        });
    }
}

function updateWishlistBadge() {
    const badge = document.getElementById('wishlist-count');
    if (badge) { badge.textContent = wishlist.length; badge.classList.toggle('show', wishlist.length > 0); }
}

// ========== BUNDLE ADD TO CART ==========
function addBundleToCart(bundleId) {
    if (bundleId === 'bundle1') { 
        addToCart(3); addToCart(4); addToCart(5); 
    } else if (bundleId === 'bundle2') { 
        addToCart(1); addToCart(2); addToCart(7); addToCart(11); 
    }
    showToast('Bundle added to cart! 🎉', 'success');
}

// ========== TESTIMONIAL SLIDER ==========
function initTestimonialSlider() {
    const track = document.getElementById('testimonial-track');
    const prevBtn = document.getElementById('testimonial-prev');
    const nextBtn = document.getElementById('testimonial-next');
    if (!track || !prevBtn || !nextBtn) return;
    let idx = 0;
    function vis() { return window.innerWidth <= 768 ? 1 : window.innerWidth <= 1024 ? 2 : 3; }
    function update() {
        const cards = track.querySelectorAll('.testimonial-card');
        const max = Math.max(0, cards.length - vis());
        if (idx > max) idx = max;
        const gap = 24;
        const w = cards[0]?.offsetWidth || 300;
        track.style.transform = `translateX(-${idx * (w + gap)}px)`;
    }
    prevBtn.addEventListener('click', () => { idx = Math.max(0, idx - 1); update(); });
    nextBtn.addEventListener('click', () => {
        const max = Math.max(0, track.querySelectorAll('.testimonial-card').length - vis());
        idx = Math.min(max, idx + 1); update();
    });
    window.addEventListener('resize', update);
    setInterval(() => {
        const max = Math.max(0, track.querySelectorAll('.testimonial-card').length - vis());
        idx = idx >= max ? 0 : idx + 1; update();
    }, 5000);
}

// ========== TOAST ==========
function showToast(message, type = '') {
    const container = document.getElementById('toast-container');
    if (!container) return;
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    toast.textContent = message;
    container.appendChild(toast);
    setTimeout(() => toast.remove(), 3000);
}

// ========== CHECKOUT FUNCTIONS ==========

// Place order function - clears cart from DATABASE after order
async function placeOrder(orderData) {
    try {
        // Check if user is logged in
        const userResult = await apiCall('/api/user', 'GET');
        if (!userResult.ok) {
            showToast('Please login to place order', 'error');
            document.getElementById('login-modal').classList.add('open');
            showLoginForm();
            return false;
        }
        
        // Add user_id and order_number to order
        orderData.user_id = userResult.data.user_id;
        orderData.order_number = 'ORD-' + Date.now().toString(36).toUpperCase();
        orderData.status = 'pending';
        orderData.order_date = new Date().toISOString();
        
        // Send order to backend
        const orderResult = await apiCall('/api/orders', 'POST', orderData);
        
        if (orderResult.ok) {
            // Cart is automatically cleared by the backend!
            // Just reload cart state to update UI
            await loadCartFromBackend();
            updateCartBadge();
            
            showToast('Order placed successfully!', 'success');
            return { success: true, order_number: orderResult.data.order_number };
        } else {
            showToast(orderResult.data.error || 'Failed to place order', 'error');
            return { success: false, error: orderResult.data.error };
        }
    } catch (error) {
        console.error('Place order error:', error);
        showToast('Network error. Please try again.', 'error');
        return { success: false, error: error.message };
    }
}

async function saveCreditCard(cardData) {
    const result = await apiCall('/api/cards', 'POST', cardData);
    if (result.ok) {
        showToast('Card saved successfully!', 'success');
        return true;
    } else {
        showToast(result.data.error || 'Failed to save card', 'error');
        return false;
    }
}

// Load user's credit cards
async function loadCreditCards() {
    const result = await apiCall('/api/cards', 'GET');
    if (result.ok && result.data) {
        return result.data;
    }
    return [];
}

// Delete credit card
async function deleteCreditCard(cardId) {
    const result = await apiCall(`/api/cards/${cardId}`, 'DELETE');
    if (result.ok) {
        showToast('Card removed', 'success');
        return true;
    } else {
        showToast('Failed to remove card', 'error');
        return false;
    }
}

// Get default credit card
async function getDefaultCreditCard() {
    const result = await apiCall('/api/cards/default', 'GET');
    if (result.ok && result.data) {
        return result.data;
    }
    return null;
}

// Credit card validation - Luhn algorithm
function validateCardNumber(cardNumber) {
    const clean = cardNumber.replace(/\s/g, '');
    if (clean.length < 13 || clean.length > 19) return false;
    
    let sum = 0;
    let alternate = false;
    for (let i = clean.length - 1; i >= 0; i--) {
        let n = parseInt(clean.charAt(i));
        if (alternate) {
            n *= 2;
            if (n > 9) n = n - 9;
        }
        sum += n;
        alternate = !alternate;
    }
    return sum % 10 === 0;
}

// Detect card brand
function detectCardBrand(cardNumber) {
    const clean = cardNumber.replace(/\s/g, '');
    const firstDigit = clean.charAt(0);
    const firstTwo = clean.substring(0, 2);
    const firstFour = clean.substring(0, 4);
    
    if (firstDigit === '4') return 'Visa';
    if (firstTwo >= '51' && firstTwo <= '55') return 'Mastercard';
    if (firstTwo === '34' || firstTwo === '37') return 'American Express';
    if (firstFour === '6011' || firstTwo === '65') return 'Discover';
    return 'Card';
}

// Format card number with spaces
function formatCardNumber(value) {
    const clean = value.replace(/\s/g, '').replace(/\D/g, '');
    let formatted = '';
    for (let i = 0; i < clean.length; i++) {
        if (i > 0 && i % 4 === 0) formatted += ' ';
        formatted += clean[i];
    }
    return formatted;
}

// Format expiry date
function formatExpiry(value) {
    const clean = value.replace(/\D/g, '');
    if (clean.length >= 2) {
        return clean.slice(0, 2) + '/' + clean.slice(2, 4);
    }
    return clean;
}

async function loadAIRecommendations() {
    try {
        const response = await fetch('/api/recommendations', {
            method: 'GET',
            credentials: 'include'
        });
        
        if (!response.ok) {
            return;
        }
        
        const products = await response.json();
        
        // Only show section if there are recommendations
        if (products && products.length > 0) {
            const section = document.getElementById('ai-recommendations-section');
            const grid = document.getElementById('ai-recommendations-grid');
            
            if (section && grid) {
                renderProductGrid('ai-recommendations-grid', products.slice(0, 4));
                section.style.display = 'block';
            }
        }
    } catch (error) {
        console.log('AI recommendations not available yet');
    }
}