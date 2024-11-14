from pydantic import BaseModel, HttpUrl, Field
from typing import List, Optional

class UserProfileBase(BaseModel):
    personal_url: Optional[HttpUrl] = None
    contact_public: Optional[bool] = None
    address: Optional[str] = None
    biography: Optional[str] = None
    organization: Optional[str] = None
    country: Optional[str] = None
    social_links: Optional[List[HttpUrl]] = None

class UserProfileCreate(UserProfileBase):
    nickname: str

class UserProfileUpdate(UserProfileBase):
    pass

class UserProfileResponse(UserProfileBase):
    id: str
    nickname: str

    class Config:
        from_attributes = True
