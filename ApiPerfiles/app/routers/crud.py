from bson import ObjectId
from pydantic_core import Url
from app.utils.database import user_profiles_collection, profile_helper
from app.schemas.user_profile import UserProfileCreate, UserProfileResponse, UserProfileUpdate
from typing import Optional

# Crear un perfil de usuario
async def create_user_profile(profile: UserProfileCreate):
    profile_dict = profile.dict()
    
    # Convertir las URL a cadenas de texto si es necesario
    if profile_dict.get("personal_url"):
        profile_dict["personal_url"] = str(profile_dict["personal_url"])
    if profile_dict.get("social_links"):
        profile_dict["social_links"] = [str(link) for link in profile_dict["social_links"]]

    result = await user_profiles_collection.insert_one(profile_dict)
    
    # Crear la respuesta con el campo `id` en lugar de `_id`
    profile_dict["id"] = str(result.inserted_id)  # Asigna el `_id` de MongoDB como `id`
    del profile_dict["_id"]  # Elimina `_id` ya que FastAPI espera `id`
    
    return profile_dict

# Obtener un perfil de usuario por nickname
async def get_user_profile(nickname: str) -> Optional[dict]:
    profile = await user_profiles_collection.find_one({"nickname": nickname})
    if profile:
        return profile_helper(profile)
    return None

# Actualizar un perfil de usuario
async def update_user_profile(nickname: str, update_data: UserProfileUpdate) -> dict:
    # Convierte el objeto de Pydantic a un diccionario plano
    update_dict = update_data.dict(exclude_unset=True)
    
    # Asegúrate de convertir cualquier campo `Url` a `str` 
    if "personal_url" in update_dict and isinstance(update_dict["personal_url"], Url):
        update_dict["personal_url"] = str(update_dict["personal_url"])
    if "social_links" in update_dict:
        update_dict["social_links"] = [str(link) for link in update_dict["social_links"]]

    # Realiza la operación de actualización en MongoDB
    result = await user_profiles_collection.update_one(
        {"nickname": nickname},
        {"$set": update_dict}
    )

    # Verifica el resultado de la actualización
    if result.modified_count == 1:
        # Retorna el perfil actualizado
        updated_profile = await user_profiles_collection.find_one({"nickname": nickname})
        # Mapea `_id` a `id`
        updated_profile['id'] = str(updated_profile.pop('_id'))
        return UserProfileResponse(**updated_profile)
    else:
        return None

# Eliminar un perfil de usuario
async def delete_user_profile(nickname: str) -> bool:
    result = await user_profiles_collection.delete_one({"nickname": nickname})
    return result.deleted_count > 0

async def get_profiles(nickname: Optional[str] = None, country: Optional[str] = None, skip: int = 0, limit: int = 10):
    query = {}
    if nickname:
        query["nickname"] = nickname
    if country:
        query["country"] = country
    profiles = user_profiles_collection.find(query).skip(skip).limit(limit)
    return [profile_helper(profile) for profile in await profiles.to_list(length=limit)]

