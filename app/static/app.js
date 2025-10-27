let currentUser = null;

async function loadUserData() {
    try {
        const response = await fetch('/api/user');
        if (!response.ok) {
            throw new Error('Not authenticated');
        }
        const data = await response.json();
        currentUser = data.user;
        updateUI(data);
    } catch (error) {
        console.error('Error loading user data:', error);
        document.getElementById('userInfo').textContent = 'Not logged in';
    }
}

function updateUI(data) {
    document.getElementById('userInfo').textContent = data.user.email;
    
    const stats = document.getElementById('usageStats');
    stats.innerHTML = `
        <div class="stat">
            <div class="stat-value">${data.usage.current_usage}</div>
            <div class="stat-label">Generated</div>
        </div>
        <div class="stat">
            <div class="stat-value">${data.usage.remaining === -1 ? '∞' : data.usage.remaining}</div>
            <div class="stat-label">Remaining</div>
        </div>
        <div class="stat">
            <div class="stat-value">${data.usage.plan_name}</div>
            <div class="stat-label">Plan</div>
        </div>
    `;
}

async function loadInvoices() {
    try {
        const response = await fetch('/api/invoices');
        if (!response.ok) return;
        
        const invoices = await response.json();
        const list = document.getElementById('invoicesList');
        
        if (!invoices || invoices.length === 0) {
            list.innerHTML = '<p class="text-muted">No invoices yet</p>';
            return;
        }
        
        list.innerHTML = invoices.map(inv => `
            <div class="invoice-item">
                <h4>${inv.invoice_number}</h4>
                <p><strong>Company:</strong> ${inv.company_name}</p>
                <p><strong>Customer:</strong> ${inv.customer_name}</p>
                <p><strong>Created:</strong> ${new Date(inv.created_at).toLocaleDateString()}</p>
            </div>
        `).join('');
    } catch (error) {
        console.error('Error loading invoices:', error);
    }
}

function addItem() {
    const container = document.getElementById('itemsContainer');
    const itemDiv = document.createElement('div');
    itemDiv.className = 'item-row';
    itemDiv.innerHTML = `
        <div class="form-row">
            <div class="form-group">
                <input type="text" class="item-name" required placeholder="Service name">
            </div>
            <div class="form-group">
                <input type="number" step="0.01" class="item-cost" required placeholder="Unit cost">
            </div>
            <div class="form-group">
                <input type="number" class="item-quantity" required value="1" placeholder="Qty">
            </div>
            <div class="form-group">
                <input type="text" class="item-description" placeholder="Description">
            </div>
        </div>
        <button type="button" class="btn-secondary" onclick="this.parentElement.remove()">Remove</button>
    `;
    container.appendChild(itemDiv);
}

document.getElementById('invoiceForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const items = [];
    document.querySelectorAll('.item-row').forEach(row => {
        items.push({
            name: row.querySelector('.item-name').value,
            description: row.querySelector('.item-description').value,
            unitCost: row.querySelector('.item-cost').value,
            quantity: row.querySelector('.item-quantity').value
        });
    });

    const data = {
        invoiceNumber: document.getElementById('invoiceNumber').value,
        companyName: document.getElementById('companyName').value,
        companyAddress: document.getElementById('companyAddress').value,
        companyCity: document.getElementById('companyCity').value,
        companyPostal: document.getElementById('companyPostal').value,
        customerName: document.getElementById('customerName').value,
        customerAddress: document.getElementById('customerAddress').value,
        customerCity: document.getElementById('customerCity').value,
        customerPostal: document.getElementById('customerPostal').value,
        customerCountry: document.getElementById('customerCountry').value,
        description: document.getElementById('description').value,
        items: items,
        notes: document.getElementById('notes').value
    };

    const result = document.getElementById('result');
    const btn = e.target.querySelector('button[type="submit"]');
    btn.disabled = true;
    btn.textContent = 'Generating...';

    try {
        const response = await fetch('/api/invoices/generate', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });

        if (response.status === 402) {
            const error = await response.json();
            result.className = 'result error';
            result.style.display = 'block';
            result.innerHTML = '<strong>Quota Exceeded:</strong> ' + error.error;
            return;
        }

        if (!response.ok) {
            throw new Error('Failed to generate invoice');
        }

        const blob = await response.blob();
        const url = window.URL.createObjectURL(blob);
        
        result.className = 'result success';
        result.style.display = 'block';
        result.innerHTML = `
            <strong>Success!</strong> Invoice generated successfully.<br>
            <a href="${url}" download="invoice.pdf" style="color: #155724; text-decoration: underline;">Download PDF</a>
        `;
        
        loadUserData();
        loadInvoices();
    } catch (error) {
        result.className = 'result error';
        result.style.display = 'block';
        result.innerHTML = '<strong>Error:</strong> ' + error.message;
    } finally {
        btn.disabled = false;
        btn.textContent = 'Generate Invoice PDF';
    }
});

async function loadPlans() {
    try {
        const response = await fetch('/api/plans');
        if (!response.ok) return;
        
        const plans = await response.json();
        const plansList = document.getElementById('plansList');
        
        plansList.innerHTML = plans.map(plan => `
            <div class="plan-card ${currentUser && currentUser.plan_id === plan.id ? 'active' : ''}">
                <div class="plan-name">${plan.name}</div>
                <div class="plan-price">$${(plan.price_cents / 100).toFixed(2)}</div>
                <div class="plan-quota">${plan.monthly_quota === -1 ? 'Unlimited' : plan.monthly_quota + ' invoices/month'}</div>
                ${currentUser && currentUser.plan_id === plan.id ? '<div><strong>Current Plan</strong></div>' : ''}
            </div>
        `).join('');
    } catch (error) {
        console.error('Error loading plans:', error);
    }
}

document.getElementById('upgradePlanBtn').addEventListener('click', () => {
    loadPlans();
    document.getElementById('plansModal').style.display = 'block';
});

function closePlansModal() {
    document.getElementById('plansModal').style.display = 'none';
}

window.onclick = function(event) {
    const modal = document.getElementById('plansModal');
    if (event.target == modal) {
        modal.style.display = 'none';
    }
}

loadUserData();
loadInvoices();
