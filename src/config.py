from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    WEATHER_API_KEY: str
    
    TG_API_KEY: str

    POSTGRES_DB: str
    POSTGRES_USER: str
    POSTGRES_PASSWORD: str
    POSTGRES_PORT: int
    POSTGRES_HOST: str

    model_config = SettingsConfigDict(env_file='/app/.env')

    @property
    def DATABASE_URL(self):
        return f"postgresql+asyncpg://{self.POSTGRES_DB}:{self.POSTGRES_PASSWORD}@{self.POSTGRES_HOST}:{self.POSTGRES_HOST}/{self.POSTGRES_DB}"

settings = Settings()