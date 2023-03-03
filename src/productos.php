 <?php include_once "includes/header.php";
    include "../conexion.php";
$id_user = $_SESSION['idUser'];
$permiso = "productos";
$sql = mysqli_query($conexion, "SELECT p.*, d.* FROM permisos p INNER JOIN detalle_permisos d ON p.id = d.id_permiso WHERE d.id_usuario = $id_user AND p.nombre = '$permiso'");
$existe = mysqli_fetch_all($sql);
if (empty($existe) && $id_user != 1) {
    header("Location: permisos.php");
}
    if (!empty($_POST)) {
		$codigo = $_POST['codigo'];
        $producto = $_POST['producto'];
        $precio = $_POST['precio'];
        $cantidad = $_POST['cantidad'];
        $usuario_id = $_SESSION['idUser'];
        $precio_bruto = $_POST['precio_bruto'];
        $marca =  $_POST['marca'];
        $alert = "";
        if (empty($codigo) || empty($producto)|| empty($marca) || empty($precio) || $precio <  0 || empty($cantidad) || $cantidad < 0 || empty($precio_bruto) || $precio_bruto < 0) {
            $alert = '<div class="alert alert-danger" role="alert">
                Todos los campos son obligatorios
              </div>';
        } else {
            $query = mysqli_query($conexion, "SELECT * FROM producto WHERE codigo = '$codigo'");
            $result = mysqli_fetch_array($query);
            if ($result > 0) {
                $alert = '<div class="alert alert-warning" role="alert">
                        El código ya existe
                    </div>';
            } else {
				$query_insert = mysqli_query($conexion,"INSERT INTO producto(codigo,descripcion,marca,precio,existencia,usuario_id,precio_bruto) values ('$codigo', '$producto','$marca','$precio','$cantidad','$usuario_id', '$precio_bruto')");
                if ($query_insert) {
                    $alert = '<div class="alert alert-success" role="alert">
                Producto Registrado
              </div>';
                } else {
                    $alert = '<div class="alert alert-danger" role="alert">
                Error al registrar el producto
              </div>';
                }
            }
        }
    }
    ?>
 <button class="btn btn-primary mb-2" type="button" data-toggle="modal" data-target="#nuevo_producto"><i class="fas fa-plus"></i></button>
 <?php echo isset($alert) ? $alert : ''; ?>
 <div class="table-responsive">
     <table class="table table-striped table-bordered" id="tbl">
         <thead class="thead-dark">
             <tr>
                 <th>#</th>
                 <th>Código</th>
                 <th>Producto</th>
                 <th>Marca</th>
                 <th>Precio</th>
                 <th>Stock</th>
                 <th>Estado</th>
                 <th>Precio Bruto</th>
                 <th></th>
             </tr>
         </thead>
         <tbody>
             <?php
                include "../conexion.php";

                $query = mysqli_query($conexion, "SELECT * FROM producto");
                $result = mysqli_num_rows($query);
                if ($result > 0) {
                    while ($data = mysqli_fetch_assoc($query)) {
                        if ($data['estado'] == 1) {
                            $estado = '<span class="badge badge-pill badge-success">Activo</span>';
                        } else {
                            $estado = '<span class="badge badge-pill badge-danger">Inactivo</span>';
                        }
                ?>
                     <tr>
                         <td><?php echo $data['codproducto']; ?></td>
                         <td><?php echo $data['codigo']; ?></td>
                         <td><?php echo $data['descripcion']; ?></td>
                         <td><?php echo $data['marca']; ?></td>
                         <td><?php echo $data['precio']; ?></td>
                         <td><?php echo $data['existencia']; ?></td>
                         <td><?php echo $estado ?></td>
                         <td><?php echo $data['precio_bruto']; ?></td>
                         <td>
                             <?php if ($data['estado'] == 1) { ?>
                                 <a href="agregar_producto.php?id=<?php echo $data['codproducto']; ?>" class="btn btn-primary"><i class='fas fa-audio-description'></i></a>

                                 <a href="editar_producto.php?id=<?php echo $data['codproducto']; ?>" class="btn btn-success"><i class='fas fa-edit'></i></a>

                                 <form action="eliminar_producto.php?id=<?php echo $data['codproducto']; ?>" method="post" class="confirmar d-inline">
                                     <button class="btn btn-danger" type="submit"><i class='fas fa-trash-alt'></i> </button>
                                 </form>
                             <?php } ?>

                             <?php if ($data['estado'] == 0) { ?>
                                 <a href="activar_producto.php?id=<?php echo $data['codproducto']; ?>" class="btn btn-success"><i class='fas fa-edit'></i></a>
                             <?php } ?>


                         </td>
                     </tr>
             <?php }
                } ?>
         </tbody>

     </table>
 </div>
 <div id="nuevo_producto" class="modal fade" tabindex="-1" role="dialog" aria-labelledby="my-modal-title" aria-hidden="true">
     <div class="modal-dialog" role="document">
         <div class="modal-content">
             <div class="modal-header bg-primary text-white">
                 <h5 class="modal-title" id="my-modal-title">Nuevo Producto</h5>
                 <button class="close" data-dismiss="modal" aria-label="Close">
                     <span aria-hidden="true">&times;</span>
                 </button>
             </div>
             <div class="modal-body">
                 <form action="" method="post" autocomplete="off">
                     <?php echo isset($alert) ? $alert : ''; ?>
                     <div class="form-group">
                         <label for="codigo">Código de Barras</label>
                         <input type="text" placeholder="Ingrese código de barras" name="codigo" id="codigo" class="form-control">
                     </div>
                     <div class="form-group">
                         <label for="producto">Producto</label>
                         <input type="text" placeholder="Ingrese nombre del producto" name="producto" id="producto" class="form-control">
                     </div>
                     <div class="form-group">
                         <label for="producto">Marca</label>
                         <input type="text" placeholder="Ingrese la marca" name="marca" id="marca" class="form-control">
                     </div>
                     <div class="form-group">
                         <label for="precio">Precio</label>
                         <input type="text" placeholder="Ingrese precio" class="form-control" name="precio" id="precio">
                     </div>
                     <div class="form-group">
                         <label for="precio_bruto">Precio Bruto</label>
                         <input type="text" placeholder="Ingrese precio Bruto" class="form-control" name="precio_bruto" id="precio_bruto">
                     </div>
                     <div class="form-group">
                         <label for="cantidad">Stock</label>
                         <input type="number" placeholder="Ingrese Stock" class="form-control" name="cantidad" id="cantidad">
                     </div>
                     <input type="submit" value="Guardar Producto" class="btn btn-primary">
                 </form>
             </div>
         </div>
     </div>
 </div>
 <div class="row">
    <div class="col-lg-12">
        <div class="form-group">
            <h4 class="text-center">Actualizar porcentaje por marca</h4><br>
        </div>
    </div>
</div>
<div id="prueba"></div>
<form method="post" id="form_marca">
                <div class="row justify-content-center">
                    <div class="col-md-2 text-center">
                        <div class="card">
                            <div class="card-header">
                                Colocar Marca
                            </div>
                            <div class="card-body">
                                <div class="form-group">
                                    <input id="idventa" class="form-control" type="text" name="id_marca" placeholder="Ingresá la marca">
                                </div>
                            </div>
                        </div>
                    </div>
                    <div class="col-md-2 text-center">
                        <div class="card">
                            <div class="card-header">
                                Colocar Porcentaje
                            </div>
                            <div class="card-body">
                                <div class="form-group">
                                    <input id="idcristal" class="form-control" type="number" name="id_porcentaje" placeholder="Ingresá el porcentaje">
                                </div>
                            </div>
                        </div>
                    </div>                  
                </div>
            </form>
            <br> 
<div class="row justify-content-center">
    <input type="button" class="btn btn-primary" value="Actualizar precio" id="btn_marca" name="btn_marca" onclick=""></input> 
</div>
<div id="prueba"></div>
 <?php include_once "includes/footer.php"; ?>
 <script >
    
    $("#btn_marca").click(function(){
    
            $.ajax({
                    url: "actualizar_porcentaje.php",
                    type: "post",
                    data: $("#form_marca").serialize(),
                    success: function(resultado){
                            $("#prueba").html(resultado);
    
                    }
    
    
            });
    
    
    
    
    
    
    });
    
    
    
    
    </script>