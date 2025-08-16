from openai import OpenAI
import chromadb

client = OpenAI(api_key=os.getenv("OPENAI_API_KEY"))

chroma_client = chromadb.Client()
collection = chroma_client.create_collection(name="with_embeddings")

text = "Chroma works well with embeddings."
embedding = client.embeddings.create(input=text, model="text-embedding-ada-002").data[0].embedding

print(embedding)

collection.add(
    documents=[text],
    embeddings=[embedding],
    ids=["embed1"]
)

query = client.embeddings.create(input="What does Chroma do?", model="text-embedding-ada-002").data[0].embedding
results = collection.query(query_embeddings=[query], n_results=1)
print(results)
