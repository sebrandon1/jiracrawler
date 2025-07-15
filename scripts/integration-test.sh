#!/bin/bash

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Script configuration
BINARY="./jiracrawler"
CONFIG_FILE=".jiracrawler-config.yaml"
TEST_USER="${JIRACRAWLER_TEST_USER:-}"
TEST_PROJECT="${JIRACRAWLER_TEST_PROJECT:-CNF}"

echo -e "${BLUE}============================================${NC}"
echo -e "${BLUE}  JiraCrawler Integration Tests${NC}"
echo -e "${BLUE}============================================${NC}"
echo ""

# Check if binary exists
if [ ! -f "$BINARY" ]; then
    echo -e "${RED}❌ Error: Binary '$BINARY' not found. Run 'make build' first.${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Binary found: $BINARY${NC}"

# Check if config file exists
if [ ! -f "$CONFIG_FILE" ]; then
    echo -e "${RED}❌ Error: Config file '$CONFIG_FILE' not found.${NC}"
    echo -e "${YELLOW}📋 To set up configuration, run:${NC}"
    echo -e "${YELLOW}   $BINARY config set --user <your-email> --token <your-api-token> --url https://issues.redhat.com${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Config file found: $CONFIG_FILE${NC}"

# Test 1: Validate credentials
echo ""
echo -e "${BLUE}🔍 Test 1: Validating JIRA credentials...${NC}"
if $BINARY validate; then
    echo -e "${GREEN}✅ Credentials validation passed${NC}"
else
    echo -e "${RED}❌ Credentials validation failed${NC}"
    echo -e "${YELLOW}📋 Please check your JIRA credentials in $CONFIG_FILE${NC}"
    exit 1
fi

# Test 2: Get assigned issues (requires test user)
echo ""
echo -e "${BLUE}🔍 Test 2: Testing assigned issues retrieval...${NC}"
if [ -z "$TEST_USER" ]; then
    echo -e "${YELLOW}⚠️  Skipping assigned issues test - JIRACRAWLER_TEST_USER not set${NC}"
    echo -e "${YELLOW}📋 To enable this test, set: export JIRACRAWLER_TEST_USER=your-email@example.com${NC}"
else
    echo -e "${BLUE}   Testing with user: $TEST_USER, project: $TEST_PROJECT${NC}"
    
    # Test JSON output
    if $BINARY get assignedissues "$TEST_USER" --projectID "$TEST_PROJECT" --output json > /tmp/jiracrawler_test_output.json; then
        echo -e "${GREEN}✅ Assigned issues JSON retrieval successful${NC}"
        
        # Basic validation of JSON output
        if command -v jq >/dev/null 2>&1; then
            if jq empty /tmp/jiracrawler_test_output.json 2>/dev/null; then
                echo -e "${GREEN}✅ JSON output is valid${NC}"
            else
                echo -e "${RED}❌ JSON output is invalid${NC}"
                cat /tmp/jiracrawler_test_output.json
                exit 1
            fi
        fi
    else
        echo -e "${RED}❌ Assigned issues retrieval failed${NC}"
        exit 1
    fi
    
    # Test YAML output
    if $BINARY get assignedissues "$TEST_USER" --projectID "$TEST_PROJECT" --output yaml > /tmp/jiracrawler_test_output.yaml; then
        echo -e "${GREEN}✅ Assigned issues YAML retrieval successful${NC}"
    else
        echo -e "${RED}❌ Assigned issues YAML retrieval failed${NC}"
        exit 1
    fi
fi

# Test 3: Test user updates in date range
echo ""
echo -e "${BLUE}🔍 Test 3: Testing user updates in date range...${NC}"
if [ -z "$TEST_USER" ]; then
    echo -e "${YELLOW}⚠️  Skipping user updates test - JIRACRAWLER_TEST_USER not set${NC}"
else
    # Test with a recent date range (last 30 days)
    START_DATE=$(date -d "30 days ago" +"%Y-%m-%d" 2>/dev/null || date -v-30d +"%Y-%m-%d" 2>/dev/null || echo "2024-01-01")
    END_DATE=$(date +"%Y-%m-%d")
    
    echo -e "${BLUE}   Testing with user: $TEST_USER, date range: $START_DATE to $END_DATE${NC}"
    
    # Test JSON output
    if $BINARY get userupdates "$TEST_USER" "$START_DATE" "$END_DATE" --output json > /tmp/jiracrawler_userupdates_test.json; then
        echo -e "${GREEN}✅ User updates JSON retrieval successful${NC}"
        
        # Basic validation of JSON output
        if command -v jq >/dev/null 2>&1; then
            if jq empty /tmp/jiracrawler_userupdates_test.json 2>/dev/null; then
                echo -e "${GREEN}✅ JSON output is valid${NC}"
                
                # Show summary
                TOTAL_COUNT=$(jq -r '.[0].totalCount // 0' /tmp/jiracrawler_userupdates_test.json 2>/dev/null || echo "0")
                echo -e "${BLUE}   Found $TOTAL_COUNT issues updated in the date range${NC}"
            else
                echo -e "${RED}❌ JSON output is invalid${NC}"
                cat /tmp/jiracrawler_userupdates_test.json
                exit 1
            fi
        fi
    else
        echo -e "${RED}❌ User updates retrieval failed${NC}"
        exit 1
    fi
    
    # Test YAML output
    if $BINARY get userupdates "$TEST_USER" "$START_DATE" "$END_DATE" --output yaml > /tmp/jiracrawler_userupdates_test.yaml; then
        echo -e "${GREEN}✅ User updates YAML retrieval successful${NC}"
    else
        echo -e "${RED}❌ User updates YAML retrieval failed${NC}"
        exit 1
    fi
fi

# Test 4: Test invalid date format handling
echo ""
echo -e "${BLUE}🔍 Test 4: Testing invalid date format handling...${NC}"
if [ -n "$TEST_USER" ]; then
    # This should fail gracefully
    if $BINARY get userupdates "$TEST_USER" "invalid-date" "2024-12-31" --output json 2>/dev/null; then
        echo -e "${RED}❌ Should have failed with invalid date format${NC}"
        exit 1
    else
        echo -e "${GREEN}✅ Invalid date format correctly rejected${NC}"
    fi
fi

# Test 5: Test help commands
echo ""
echo -e "${BLUE}🔍 Test 5: Testing help commands...${NC}"
if $BINARY --help > /dev/null; then
    echo -e "${GREEN}✅ Main help command works${NC}"
else
    echo -e "${RED}❌ Main help command failed${NC}"
    exit 1
fi

if $BINARY get --help > /dev/null; then
    echo -e "${GREEN}✅ Get help command works${NC}"
else
    echo -e "${RED}❌ Get help command failed${NC}"
    exit 1
fi

if $BINARY get userupdates --help > /dev/null; then
    echo -e "${GREEN}✅ User updates help command works${NC}"
else
    echo -e "${RED}❌ User updates help command failed${NC}"
    exit 1
fi

# Cleanup
echo ""
echo -e "${BLUE}🧹 Cleaning up temporary files...${NC}"
rm -f /tmp/jiracrawler_test_output.json /tmp/jiracrawler_test_output.yaml
rm -f /tmp/jiracrawler_userupdates_test.json /tmp/jiracrawler_userupdates_test.yaml

# Success
echo ""
echo -e "${GREEN}============================================${NC}"
echo -e "${GREEN}🎉 All integration tests passed!${NC}"
echo -e "${GREEN}============================================${NC}"
echo ""

if [ -z "$TEST_USER" ]; then
    echo -e "${YELLOW}💡 Tip: Set JIRACRAWLER_TEST_USER environment variable to enable full testing${NC}"
    echo -e "${YELLOW}   Example: export JIRACRAWLER_TEST_USER=your-email@example.com${NC}"
fi

echo -e "${BLUE}📊 Integration test summary:${NC}"
echo -e "${BLUE}  - Credentials validation: ✅${NC}"
echo -e "${BLUE}  - Assigned issues retrieval: $([ -n "$TEST_USER" ] && echo "✅" || echo "⚠️ (skipped)")${NC}"
echo -e "${BLUE}  - User updates retrieval: $([ -n "$TEST_USER" ] && echo "✅" || echo "⚠️ (skipped)")${NC}"
echo -e "${BLUE}  - Error handling: ✅${NC}"
echo -e "${BLUE}  - Help commands: ✅${NC}"
echo "" 