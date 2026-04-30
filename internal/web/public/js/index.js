var id = 0;
var today = null;

function addExercise(name, weight, reps, intensity, color) {
    id++;
    var container = document.getElementById("todayEx");
    var entry = document.createElement('div');
    entry.className = 'workout-entry';
    entry.id = 'entry-' + id;

    var safeColor = color || '#6c757d';
    var safeIntensity = (intensity !== undefined && intensity !== '') ? intensity : 5;

    entry.innerHTML = `
        <div class="entry-color-strip" style="background-color:${safeColor};"></div>
        <input type="hidden" name="name" value="${name}">
        <input type="hidden" name="weight" value="${weight || 0}">
        <input type="hidden" name="reps" value="${reps || 0}">
        <span class="entry-name" title="${name}">${name}</span>
        <div class="entry-controls">
            <span class="entry-label">Intensity</span>
            <input type="number" class="form-control entry-intensity"
                name="intensity" value="${safeIntensity}" min="0" max="10">
            <input type="color" class="entry-color-picker"
                name="workout_color" value="${safeColor}"
                oninput="this.previousElementSibling.previousElementSibling.previousElementSibling.previousElementSibling.style.backgroundColor=this.value">
            <button type="button" class="entry-del-btn"
                onclick="document.getElementById('entry-${id}').remove(); updateEmptyState();"
                title="Remove">
                <i class="bi bi-x-lg"></i>
            </button>
        </div>
    `;

    // Wire up color picker → color strip live update
    var strip = entry.querySelector('.entry-color-strip');
    var colorPicker = entry.querySelector('.entry-color-picker');
    colorPicker.addEventListener('input', function() {
        strip.style.backgroundColor = this.value;
    });
    // Remove the inline oninput we set above (cleaner)
    colorPicker.removeAttribute('oninput');

    container.appendChild(entry);
    updateEmptyState();
}

function updateEmptyState() {
    var hasEntries = document.getElementById('todayEx').children.length > 0;
    document.getElementById('emptyState').style.display = hasEntries ? 'none' : '';
}

function setFormContent(sets, date) {
    window.sessionStorage.setItem("today", date);
    today = date;
    document.getElementById('todayEx').innerHTML = "";
    updateEmptyState();
    document.getElementById("formDate").value = date;
    document.getElementById("realDate").value = date;

    // Update heatmap highlights
    if (window.intensityChart) {
        window.intensityChart.data.datasets[0].selectedDate = date;
        window.intensityChart.update();
    }
    if (window.colorChart) {
        window.colorChart.data.datasets[0].selectedDate = date;
        window.colorChart.update();
    }

    if (sets) {
        for (var i = 0; i < sets.length; i++) {
            if (sets[i].Date == date) {
                addExercise(sets[i].Name, sets[i].Weight, sets[i].Reps, sets[i].Intensity, sets[i].WorkoutColor);
            }
        }
    }
}

function setFormDate(sets) {
    var date = window.sessionStorage.getItem("today");
    if (!date) {
        date = new Date().toISOString().split('T')[0];
    }
    setFormContent(sets, date);
}

function goToToday() {
    var date = new Date().toISOString().split('T')[0];
    setFormContent(window._allSets, date);
}

function setWeightDate() {
    var date = document.getElementById("realDate").value;
    document.getElementById("weightDate").value = date;
}

function delExercise(exID) {
    document.getElementById(exID).remove();
    updateEmptyState();
}

function moveDayLeftRight(where, sets) {
    var dateStr = document.getElementById("realDate").value;
    var year  = dateStr.substring(0, 4);
    var month = dateStr.substring(5, 7);
    var day   = dateStr.substring(8, 10);
    var date  = new Date(year, month - 1, day);
    date.setDate(date.getDate() + parseInt(where));
    var newDate = date.toLocaleDateString('en-CA');
    setFormContent(sets, newDate);
}

function addAllGroup(exs, gr) {
    if (!exs) return;
    for (var i = 0; i < exs.length; i++) {
        if (exs[i].Group == gr) {
            addExercise(exs[i].Name, exs[i].Weight, exs[i].Reps, exs[i].Intensity, exs[i].Color);
        }
    }
}

function selectGroup(gr) {
    window._selectedGroup = gr;

    // Highlight active chip
    document.querySelectorAll('.group-chip').forEach(function(chip) {
        chip.classList.toggle('active', chip.getAttribute('data-group') === gr);
    });

    // Show group header, hide "no group" state
    var header = document.getElementById('groupHeader');
    var noGroup = document.getElementById('noGroupState');
    var searchInput = document.getElementById('exSearch');
    if (header) {
        header.style.display = 'flex';
        document.getElementById('groupHeaderName').textContent = gr;
    }
    if (noGroup) noGroup.style.display = 'none';
    if (searchInput) searchInput.value = '';

    // Show only items in this group
    document.querySelectorAll('.exercise-item').forEach(function(item) {
        item.style.display = item.getAttribute('data-group') === gr ? '' : 'none';
    });
}

function clearGroup() {
    window._selectedGroup = null;

    // Deactivate all chips
    document.querySelectorAll('.group-chip').forEach(function(chip) {
        chip.classList.remove('active');
    });

    // Hide group header, show "no group" state, hide all exercises
    var header = document.getElementById('groupHeader');
    var noGroup = document.getElementById('noGroupState');
    if (header) header.style.display = 'none';
    if (noGroup) noGroup.style.display = '';
    document.querySelectorAll('.exercise-item').forEach(function(item) {
        item.style.display = 'none';
    });
}

function filterExercises() {
    var query = document.getElementById('exSearch').value.toLowerCase().trim();
    var gr = window._selectedGroup;
    if (!gr) return;

    document.querySelectorAll('.exercise-item').forEach(function(item) {
        var nameMatch = item.getAttribute('data-name').toLowerCase().includes(query);
        var groupMatch = item.getAttribute('data-group') === gr;
        item.style.display = (groupMatch && (!query || nameMatch)) ? '' : 'none';
    });
}
