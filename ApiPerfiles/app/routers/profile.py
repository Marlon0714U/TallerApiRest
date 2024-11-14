from typing import List, Optional
from fastapi import APIRouter, HTTPException, Query
from app.schemas.user_profile import UserProfileCreate, UserProfileResponse, UserProfileUpdate
from app.routers import crud

router = APIRouter()

@router.post("/", response_model=UserProfileResponse)
async def create_profile(profile: UserProfileCreate):
    existing_profile = await crud.get_user_profile(profile.nickname)
    if existing_profile:
        raise HTTPException(status_code=400, detail="Nickname already exists")
    return await crud.create_user_profile(profile)

@router.get("/", response_model=List[UserProfileResponse])
async def get_profiles(
    nickname: Optional[str] = None,
    country: Optional[str] = None,
    skip: int = Query(0, ge=0),
    limit: int = Query(10, ge=1)
):
    profiles = await crud.get_profiles(nickname=nickname, country=country, skip=skip, limit=limit)
    return profiles

@router.get("/{nickname}", response_model=UserProfileResponse)
async def get_profile(nickname: str):
    profile = await crud.get_user_profile(nickname)
    if not profile:
        raise HTTPException(status_code=404, detail="Profile not found")
    return profile

@router.put("/{nickname}", response_model=UserProfileResponse)
async def update_profile(nickname: str, profile_update: UserProfileUpdate):
    profile = await crud.update_user_profile(nickname, profile_update)
    if not profile:
        raise HTTPException(status_code=404, detail="Profile not found")
    return profile

@router.delete("/{nickname}", response_model=dict)
async def delete_profile(nickname: str):
    deleted = await crud.delete_user_profile(nickname)
    if not deleted:
        raise HTTPException(status_code=404, detail="Profile not found")
    return {"message": "Profile deleted successfully"}
