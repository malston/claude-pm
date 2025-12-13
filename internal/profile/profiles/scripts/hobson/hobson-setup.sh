#!/bin/bash
# ABOUTME: Interactive setup wizard for the hobson profile
# ABOUTME: Guides users through category-based plugin selection for wshobson/agents

set -e

# Check if gum is available for interactive UI
HAS_GUM=false
if command -v gum &> /dev/null; then
    HAS_GUM=true
fi

# Marketplace suffix for plugins
MARKETPLACE="wshobson-agents"

# Category definitions - maps theme names to plugin lists
declare -A CATEGORIES
CATEGORIES["Core Development"]="code-documentation debugging-toolkit git-pr-workflows backend-development frontend-mobile-development full-stack-orchestration code-refactoring dependency-management error-debugging team-collaboration documentation-generation c4-architecture multi-platform-apps developer-essentials"
CATEGORIES["Quality & Testing"]="unit-testing tdd-workflows code-review-ai comprehensive-review performance-testing-review framework-migration codebase-cleanup"
CATEGORIES["AI & Machine Learning"]="llm-application-dev agent-orchestration context-management machine-learning-ops"
CATEGORIES["Infrastructure & DevOps"]="deployment-strategies deployment-validation kubernetes-operations cloud-infrastructure cicd-automation incident-response error-diagnostics distributed-debugging observability-monitoring"
CATEGORIES["Security & Compliance"]="security-scanning security-compliance backend-api-security frontend-mobile-security"
CATEGORIES["Data & Databases"]="data-engineering data-validation-suite database-design database-migrations application-performance database-cloud-optimization"
CATEGORIES["Languages"]="python-development javascript-typescript systems-programming jvm-languages web-scripting functional-programming julia-development arm-cortex-microcontrollers shell-scripting"
CATEGORIES["Business & Specialty"]="api-scaffolding api-testing-observability seo-content-creation seo-technical-optimization seo-analysis-monitoring business-analytics hr-legal-compliance customer-sales-automation content-marketing blockchain-web3 quantitative-trading payment-processing game-development accessibility-compliance"

# Category descriptions for display
declare -A CATEGORY_DESCRIPTIONS
CATEGORY_DESCRIPTIONS["Core Development"]="workflows, debugging, docs, refactoring"
CATEGORY_DESCRIPTIONS["Quality & Testing"]="code review, testing, cleanup"
CATEGORY_DESCRIPTIONS["AI & Machine Learning"]="LLM dev, agents, MLOps"
CATEGORY_DESCRIPTIONS["Infrastructure & DevOps"]="K8s, cloud, CI/CD, monitoring"
CATEGORY_DESCRIPTIONS["Security & Compliance"]="scanning, compliance, API security"
CATEGORY_DESCRIPTIONS["Data & Databases"]="ETL, schema design, migrations"
CATEGORY_DESCRIPTIONS["Languages"]="Python, JS/TS, Go, Rust, etc."
CATEGORY_DESCRIPTIONS["Business & Specialty"]="SEO, analytics, blockchain, gaming"

# Selected categories (space-separated names)
SELECTED_CATEGORIES=""

# Print header
print_header() {
    echo ""
    echo "╭─────────────────────────────────────────────────────╮"
    echo "│       Welcome to the Hobson Profile Setup           │"
    echo "╰─────────────────────────────────────────────────────╯"
    echo ""
}

# Print category list with checkboxes
print_categories() {
    local idx=1
    for category in "Core Development" "Quality & Testing" "AI & Machine Learning" "Infrastructure & DevOps" "Security & Compliance" "Data & Databases" "Languages" "Business & Specialty"; do
        local marker="[ ]"
        if [[ " $SELECTED_CATEGORIES " == *" $category "* ]]; then
            marker="[x]"
        fi
        local desc="${CATEGORY_DESCRIPTIONS[$category]}"
        printf "  %d. %s %-25s (%s)\n" "$idx" "$marker" "$category" "$desc"
        ((idx++))
    done
}

# Toggle category selection
toggle_category() {
    local category="$1"
    if [[ " $SELECTED_CATEGORIES " == *" $category "* ]]; then
        # Remove category
        SELECTED_CATEGORIES="${SELECTED_CATEGORIES// $category / }"
        SELECTED_CATEGORIES="${SELECTED_CATEGORIES/#$category /}"
        SELECTED_CATEGORIES="${SELECTED_CATEGORIES/% $category/}"
        SELECTED_CATEGORIES="${SELECTED_CATEGORIES/$category/}"
    else
        # Add category
        if [[ -z "$SELECTED_CATEGORIES" ]]; then
            SELECTED_CATEGORIES="$category"
        else
            SELECTED_CATEGORIES="$SELECTED_CATEGORIES $category"
        fi
    fi
}

# Interactive selection using gum
select_with_gum() {
    local categories=("Core Development" "Quality & Testing" "AI & Machine Learning" "Infrastructure & DevOps" "Security & Compliance" "Data & Databases" "Languages" "Business & Specialty")
    local options=()

    for category in "${categories[@]}"; do
        local desc="${CATEGORY_DESCRIPTIONS[$category]}"
        options+=("$category ($desc)")
    done

    echo "Select development categories (space to toggle, enter to confirm):"
    echo ""

    local selected
    selected=$(gum choose --no-limit "${options[@]}" 2>/dev/null || true)

    if [[ -z "$selected" ]]; then
        return 1
    fi

    # Parse selected options back to category names
    SELECTED_CATEGORIES=""
    while IFS= read -r line; do
        # Extract category name before the parenthesis
        local cat_name="${line%% (*}"
        if [[ -n "$cat_name" ]]; then
            if [[ -z "$SELECTED_CATEGORIES" ]]; then
                SELECTED_CATEGORIES="$cat_name"
            else
                SELECTED_CATEGORIES="$SELECTED_CATEGORIES|$cat_name"
            fi
        fi
    done <<< "$selected"

    return 0
}

# Interactive selection using basic prompts (fallback)
select_with_prompts() {
    local categories=("Core Development" "Quality & Testing" "AI & Machine Learning" "Infrastructure & DevOps" "Security & Compliance" "Data & Databases" "Languages" "Business & Specialty")

    while true; do
        clear
        print_header
        echo "Select development categories (enter numbers to toggle):"
        echo ""
        print_categories
        echo ""
        echo "Commands: [1-8] Toggle category  [a] Select all  [n] Select none"
        echo "          [c] Customize plugins  [Enter] Continue  [q] Quit"
        echo ""
        read -p "> " input

        case "$input" in
            [1-8])
                local idx=$((input - 1))
                toggle_category "${categories[$idx]}"
                ;;
            a|A)
                SELECTED_CATEGORIES="Core Development Quality & Testing AI & Machine Learning Infrastructure & DevOps Security & Compliance Data & Databases Languages Business & Specialty"
                ;;
            n|N)
                SELECTED_CATEGORIES=""
                ;;
            c|C)
                if [[ -z "$SELECTED_CATEGORIES" ]]; then
                    echo "Please select at least one category first."
                    sleep 1
                else
                    return 2  # Signal to customize
                fi
                ;;
            q|Q)
                echo "Setup cancelled."
                exit 0
                ;;
            "")
                if [[ -z "$SELECTED_CATEGORIES" ]]; then
                    echo "Please select at least one category."
                    sleep 1
                else
                    return 0  # Continue without customizing
                fi
                ;;
        esac
    done
}

# Enable plugins for selected categories
enable_plugins() {
    local plugins_to_enable=()

    # Collect all plugins from selected categories
    # Handle both space-separated (basic) and pipe-separated (gum) formats
    local IFS_ORIG="$IFS"
    if [[ "$SELECTED_CATEGORIES" == *"|"* ]]; then
        IFS="|"
    else
        # For space-separated, we need to handle multi-word category names
        # Convert to array manually
        local temp_categories=()
        for category in "Core Development" "Quality & Testing" "AI & Machine Learning" "Infrastructure & DevOps" "Security & Compliance" "Data & Databases" "Languages" "Business & Specialty"; do
            if [[ " $SELECTED_CATEGORIES " == *"$category"* ]]; then
                temp_categories+=("$category")
            fi
        done

        for category in "${temp_categories[@]}"; do
            local category_plugins="${CATEGORIES[$category]}"
            for plugin in $category_plugins; do
                plugins_to_enable+=("$plugin")
            done
        done

        IFS="$IFS_ORIG"

        # Enable the plugins
        echo ""
        echo "Enabling ${#plugins_to_enable[@]} plugins..."
        echo ""

        local success=0
        local failed=0

        for plugin in "${plugins_to_enable[@]}"; do
            local full_name="${plugin}@${MARKETPLACE}"
            if claudeup enable "$full_name" > /dev/null 2>&1; then
                echo "  ✓ $full_name"
                ((success++))
            else
                echo "  ✗ $full_name (failed)"
                ((failed++))
            fi
        done

        echo ""
        echo "Setup complete!"
        echo "  Enabled: $success plugins"
        if [[ $failed -gt 0 ]]; then
            echo "  Failed:  $failed plugins"
        fi
        echo ""
        echo "Run 'claudeup status' to see your configuration."
        return
    fi

    # For pipe-separated (from gum)
    for category in $SELECTED_CATEGORIES; do
        local category_plugins="${CATEGORIES[$category]}"
        for plugin in $category_plugins; do
            plugins_to_enable+=("$plugin")
        done
    done
    IFS="$IFS_ORIG"

    # Enable the plugins
    echo ""
    echo "Enabling ${#plugins_to_enable[@]} plugins..."
    echo ""

    local success=0
    local failed=0

    for plugin in "${plugins_to_enable[@]}"; do
        local full_name="${plugin}@${MARKETPLACE}"
        if claudeup enable "$full_name" > /dev/null 2>&1; then
            echo "  ✓ $full_name"
            ((success++))
        else
            echo "  ✗ $full_name (failed)"
            ((failed++))
        fi
    done

    echo ""
    echo "Setup complete!"
    echo "  Enabled: $success plugins"
    if [[ $failed -gt 0 ]]; then
        echo "  Failed:  $failed plugins"
    fi
    echo ""
    echo "Run 'claudeup status' to see your configuration."
}

# Main execution
main() {
    print_header

    if [[ "$HAS_GUM" == "true" ]]; then
        if select_with_gum; then
            enable_plugins
        else
            echo "No categories selected. Setup cancelled."
            exit 0
        fi
    else
        local result
        select_with_prompts
        result=$?

        if [[ $result -eq 0 ]]; then
            enable_plugins
        elif [[ $result -eq 2 ]]; then
            # TODO: Implement per-plugin customization
            echo "Per-plugin customization coming soon..."
            enable_plugins
        fi
    fi
}

main "$@"
