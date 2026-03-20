# General
All database entities and data transfer objects  with some nearly depended business logic functions

## programming rules
- All models who are store within the database should have postfix "Entity".
- All models that are used as DTO should have the postifx "Dto"
- A Collection root model like License/Project/User should have domain.RootEntity implemented
- A Collection children model like ProjectVersion should have domain.ChildEntity implemented when we need a unique id
 
