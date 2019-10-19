# item-API
# Golang Workshop - Challenge
## Day 3  
### Ejercicio

> Usando la documentación para desarrolladores de mercadolibre:
> https://developers.mercadolibre.com.ar/  
> Dado un item, necesitamos visualizar información sobre el mismo, su vendedor, categoría y site. Para eso proponemos crear una REST Api.

>Input: item_id  
>Output: Un JSON con la información completa del item, incluyendo la información relativa al site, vendedor y categoría correspondientes.
>Además, necesitamos poder pasar opcionalmente un parámetro “attributes” con una lista de campos, para obtener sólo una parte de la respuesta ampliada. 
> (Valores posibles: “site”, “seller” y/o “category”)  
>Endpoints de utilidad:  
https://api.mercadolibre.com/items/ :item_id  
https://api.mercadolibre.com/sites/: site_id  
https://api.mercadolibre.com/users/: user_id  
https://api.mercadolibre.com/categories/ :category_id  

#### Start:
```sh
$ go run .
```
#
## item-API:
### Ejemplos:
##### test service
```sh
GET: localhost:8080/ping
```
##### respuesta full
```sh
GET: localhost:8080/show/MLA726125948
```
##### con querystring (respuesta solo con seller y category)
```sh
GET: localhost:8080/show/ MLA726125948 ?attributes=seller,category
```
