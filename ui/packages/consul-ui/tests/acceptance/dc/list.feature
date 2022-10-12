@setupApplicationTest
# Nodes have been removed from this list because they randomly get
# assigned the synthetic-node property which filters them out of the list page.
# We are already testing the nodes being listed in the nodes index feature test.
Feature: dc / list: List Models
  Scenario: Listing [Model]
    Given 1 datacenter model with the value "dc-1"
    And 3 [Model] models
    When I visit the [Page] page for yaml
    ---
      dc: dc-1
    ---
    Then the url should be [Url]

    Then I see 3 [Model] models
  Where:
    -------------------------------------------------
    | Model   | Page     | Url                      |
    # | node    | nodes    | /dc-1/nodes              |
    | kv      | kvs      | /dc-1/kv                 |
    # | acl     | acls     | /dc-1/acls               |
    | token   | tokens   | /dc-1/acls/tokens        |
    | policy  | policies | /dc-1/acls/policies      |
    -------------------------------------------------
