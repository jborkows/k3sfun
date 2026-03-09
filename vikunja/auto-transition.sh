#!/bin/bash
set -euo pipefail

# Vikunja Auto-Transition Sidecar
# Automatically transitions tasks based on:
# 1. Blocked tasks - moves to ready when blockers are completed
# 2. Scheduled tasks - moves to ready when earliest-on date is reached

VIKUNJA_API_URL="${VIKUNJA_API_URL:-http://localhost:3456}"
VIKUNJA_API_TOKEN="${VIKUNJA_API_TOKEN:-}"
CHECK_INTERVAL="${CHECK_INTERVAL:-300}"  # 5 minutes in seconds

# State labels (used for filtering and updating tasks)
LABEL_BLOCKED="state:blocked"
LABEL_SCHEDULED="state:scheduled"

# Logging function
log() {
  echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1"
}

# Check if API token is configured
check_config() {
  if [[ -z "$VIKUNJA_API_TOKEN" ]]; then
    log "WARNING: VIKUNJA_API_TOKEN not set. Auto-transition will be disabled."
    log "To enable, create a Vikunja API token and set the secret."
    return 1
  fi
  return 0
}

# Get current date in YYYY-MM-DD format
get_today() {
  date +%Y-%m-%d
}

# Get all tasks with a specific label
get_tasks_with_label() {
  local label="$1"
  # Query tasks by label using Vikunja API
  # Note: Vikunja API v1 uses filter syntax
  curl -s -H "Authorization: Bearer $VIKUNJA_API_TOKEN" \
    "${VIKUNJA_API_URL}/api/v1/tasks?filter_labels=${label}" 2>/dev/null || echo "[]"
}

# Get task details including relations
get_task_details() {
  local task_id="$1"
  curl -s -H "Authorization: Bearer $VIKUNJA_API_TOKEN" \
    "${VIKUNJA_API_URL}/api/v1/tasks/${task_id}" 2>/dev/null
}

# Get task relations
get_task_relations() {
  local task_id="$1"
  curl -s -H "Authorization: Bearer $VIKUNJA_API_TOKEN" \
    "${VIKUNJA_API_URL}/api/v1/tasks/${task_id}/relations" 2>/dev/null || echo "[]"
}

# Update task labels
update_task_labels() {
  local task_id="$1"
  local labels="$2"
  
  # Update labels via API
  curl -s -X PUT \
    -H "Authorization: Bearer $VIKUNJA_API_TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"labels\": ${labels}}" \
    "${VIKUNJA_API_URL}/api/v1/tasks/${task_id}" 2>/dev/null
}

# Check if task has any incomplete blockers
has_incomplete_blockers() {
  local task_id="$1"
  
  # Get all relations for this task
  local relations
  relations=$(get_task_relations "$task_id")
  
  # Check for "blocked" relations where the blocking task is not done
  # Vikunja API returns relations with "kind" field
  # "blocked" means this task is blocked by another
  echo "$relations" | jq -r '.[] | select(.kind == "blocked") | .other_task_id' 2>/dev/null | while read -r blocker_id; do
    if [[ -n "$blocker_id" ]]; then
      local blocker_data
      blocker_data=$(get_task_details "$blocker_id")
      local is_done
      is_done=$(echo "$blocker_data" | jq -r '.done // false')
      if [[ "$is_done" != "true" ]]; then
        echo "yes"
        return
      fi
    fi
  done
  echo "no"
}

# Process blocked tasks - check if blockers are resolved
process_blocked_tasks() {
  log "Checking blocked tasks..."
  
  local tasks
  tasks=$(get_tasks_with_label "$LABEL_BLOCKED")
  
  echo "$tasks" | jq -c '.[]?' 2>/dev/null | while read -r task; do
    if [[ -z "$task" ]]; then
      continue
    fi
    
    local task_id
    task_id=$(echo "$task" | jq -r '.id')
    local task_title
    task_title=$(echo "$task" | jq -r '.title')
    
    log "Checking blocked task: $task_title (ID: $task_id)"
    
    # Check if there are incomplete blockers
    local has_blockers
    has_blockers=$(has_incomplete_blockers "$task_id")
    
    if [[ "$has_blockers" == "no" ]]; then
      log "  -> No incomplete blockers found, transitioning to ready"
      
      # Get current labels
      local current_labels
      current_labels=$(echo "$task" | jq -r '.labels // []')
      
      # Remove blocked label, add ready label
      local new_labels
      new_labels=$(echo "$current_labels" | jq 'map(select(. != "state:blocked")) | . + ["state:ready"] | unique')
      
      # Update task
      if update_task_labels "$task_id" "$new_labels"; then
        log "  -> Successfully transitioned to ready"
      else
        log "  -> Failed to update task labels"
      fi
    else
      log "  -> Still has incomplete blockers"
    fi
  done
}

# Process scheduled tasks - check if earliest-on date has passed
process_scheduled_tasks() {
  log "Checking scheduled tasks..."
  
  local today
  today=$(get_today)
  
  local tasks
  tasks=$(get_tasks_with_label "$LABEL_SCHEDULED")
  
  echo "$tasks" | jq -c '.[]?' 2>/dev/null | while read -r task; do
    if [[ -z "$task" ]]; then
      continue
    fi
    
    local task_id
    task_id=$(echo "$task" | jq -r '.id')
    local task_title
    task_title=$(echo "$task" | jq -r '.title')
    
    # Get labels to find earliest-on date
    local labels
    labels=$(echo "$task" | jq -r '.labels // []')
    
    # Find earliest-on label
    local earliest_on
    earliest_on=$(echo "$labels" | jq -r '.[] | select(startswith("earliest-on:"))' | head -1)
    
    if [[ -n "$earliest_on" ]]; then
      # Extract date from label (format: earliest-on:YYYY-MM-DD)
      local scheduled_date
      scheduled_date=${earliest_on#earliest-on:}
      
      log "Task: $task_title (ID: $task_id) - scheduled for: $scheduled_date"
      
      # Compare dates (YYYY-MM-DD format allows string comparison)
      if [[ "$today" > "$scheduled_date" ]] || [[ "$today" == "$scheduled_date" ]]; then
        log "  -> Scheduled date reached, transitioning to ready"
        
        # Get current labels
        local current_labels
        current_labels=$(echo "$task" | jq -r '.labels // []')
        
        # Remove scheduled and earliest-on labels, add ready label
        local new_labels
        new_labels=$(echo "$current_labels" | jq 'map(select(. != "state:scheduled" and (startswith("earliest-on:") | not))) | . + ["state:ready"] | unique')
        
        # Update task
        if update_task_labels "$task_id" "$new_labels"; then
          log "  -> Successfully transitioned to ready"
        else
          log "  -> Failed to update task labels"
        fi
      else
        log "  -> Still waiting for $scheduled_date (today: $today)"
      fi
    fi
  done
}

# Main loop
main() {
  log "Vikunja Auto-Transition Sidecar started"
  log "API URL: $VIKUNJA_API_URL"
  log "Check interval: ${CHECK_INTERVAL}s"
  
  if ! check_config; then
    log "Waiting for API token configuration..."
    # Keep running but don't process - allows secret to be added later
    while true; do
      sleep 60
      if check_config; then
        break
      fi
    done
  fi
  
  log "Configuration valid. Starting task processing..."
  
  # Run once immediately
  process_blocked_tasks
  process_scheduled_tasks
  
  # Then run periodically
  while true; do
    log "Sleeping for ${CHECK_INTERVAL}s..."
    sleep "$CHECK_INTERVAL"
    
    log "Running task checks..."
    process_blocked_tasks
    process_scheduled_tasks
    log "Task checks complete"
  done
}

# Handle signals gracefully
trap 'log "Received signal, shutting down..."; exit 0' TERM INT

main "$@"
