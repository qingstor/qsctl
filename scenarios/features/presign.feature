Feature: Generate a link with time limit for a QingStor object.

	Scenario: Generate a link for a QingStor object with non-public authority.
		Given a non-public bucket with objects
		    | name                    | expire_seconds |
			| test_non_public_presign | 3600           |
		When generate the link with signature
		Then the link should include parameters like signature, etc.

	Scenario: Generate a link for a QingStor object with full permissions.
		Given a public bucket with objects
   		    | name                |
			| test_public_presign |
		When generate the spliced link
		Then the link should include object_key
