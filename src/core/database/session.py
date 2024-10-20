from core.database.metadata import session_factory

async def get_session():
    with session_factory() as session:
        yield session 